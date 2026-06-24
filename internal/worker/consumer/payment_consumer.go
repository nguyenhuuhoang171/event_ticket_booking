package consumer

import (
	"context"
	"log"
	"time"

	"event_ticket_booking/config"
	"event_ticket_booking/constant"
	"event_ticket_booking/internal/domain/booking/event"
	"event_ticket_booking/internal/worker/payment"
	commonModel "event_ticket_booking/model"

	"github.com/IBM/sarama"
)

// paymentConsumer xử lý message thanh toán từ topic payment.request.
type paymentConsumer struct {
	processor payment.Processor
}

func NewPaymentConsumer(cfg config.Config, lib commonModel.Lib) Subscriber {
	return paymentConsumer{processor: payment.NewProcessor(cfg, lib)}
}

func (paymentConsumer) Group() string    { return constant.KAFKA_GROUP_PAYMENT }
func (paymentConsumer) Topics() []string { return []string{constant.TOPIC_PAYMENT_REQUEST} }

func (paymentConsumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (paymentConsumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c paymentConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case m, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			msg, err := event.ParsePaymentMessage(m.Value)
			if err != nil {
				log.Printf("[Consumer] invalid payload: %v", err)
				session.MarkMessage(m, "") // invalid payload -> bỏ qua, không retry
				continue
			}

			kafkaCfg := c.processor.KafkaConfig()
			err = withRetry(kafkaCfg, func() error {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
				defer cancel()
				return c.processor.ProcessPayment(ctx, msg.BookingID)
			})
			if err != nil {
				log.Printf("[Consumer] ProcessPayment booking %d failed after %d retries: %v",
					msg.BookingID, kafkaCfg.Consumer.RetryMax, err)
			}
			session.MarkMessage(m, "")

		case <-session.Context().Done():
			return nil
		}
	}
}
