package consumer

import (
	"context"
	"fmt"
	"log"
	"time"

	"event_ticket_booking/config"
	kafkaInfra "event_ticket_booking/infrastructure/kafka"

	"github.com/IBM/sarama"
)

// withRetry gọi fn tối đa RetryMax+1 lần, nghỉ RetryBackoffMs giữa các lần thử.
func withRetry(kafkaCfg config.KafkaConfig, fn func() error) error {
	maxRetries := kafkaCfg.Consumer.RetryMax
	backoff := time.Duration(kafkaCfg.Consumer.RetryBackoffMs) * time.Millisecond

	var err error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}
		if attempt < maxRetries {
			time.Sleep(backoff)
		}
	}
	return err
}

type Subscriber interface {
	Group() string
	Topics() []string
	sarama.ConsumerGroupHandler
}

type Runner struct {
	cfg         config.KafkaConfig
	subscribers []Subscriber
	groups      []sarama.ConsumerGroup
}

func NewRunner(cfg config.KafkaConfig, subscribers ...Subscriber) *Runner {
	return &Runner{cfg: cfg, subscribers: subscribers}
}

// Start tạo group + chạy consume loop cho mỗi subscriber. Dừng khi ctx bị cancel.
func (r *Runner) Start(ctx context.Context) error {
	for _, sub := range r.subscribers {
		group, err := kafkaInfra.NewConsumerGroup(r.cfg, sub.Group())
		if err != nil {
			return fmt.Errorf("create consumer group %q: %w", sub.Group(), err)
		}
		r.groups = append(r.groups, group)

		go r.consume(ctx, group, sub)
		go r.logErrors(ctx, group, sub)
		log.Printf("[Consumer] %q subscribed to %v", sub.Group(), sub.Topics())
	}
	return nil
}

func (r *Runner) Close() error {
	for _, group := range r.groups {
		if err := group.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (r *Runner) consume(ctx context.Context, group sarama.ConsumerGroup, sub Subscriber) {
	for {
		// Consume trả về khi rebalance hoặc ctx bị cancel -> lặp lại để tiếp tục.
		if err := group.Consume(ctx, sub.Topics(), sub); err != nil {
			if ctx.Err() != nil {
				break
			}
			log.Printf("[Consumer] %q consume: %v", sub.Group(), err)
		}
		if ctx.Err() != nil {
			break
		}
	}
	log.Printf("[Consumer] %q stopped", sub.Group())
}

// logErrors drain channel lỗi bất đồng bộ (do Consumer.Return.Errors = true).
func (r *Runner) logErrors(ctx context.Context, group sarama.ConsumerGroup, sub Subscriber) {
	for {
		select {
		case <-ctx.Done():
			return
		case err, ok := <-group.Errors():
			if !ok {
				return
			}
			log.Printf("[Consumer] %q group error: %v", sub.Group(), err)
		}
	}
}
