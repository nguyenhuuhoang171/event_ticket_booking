package kafka

import (
	"errors"

	"event_ticket_booking/config"
	"event_ticket_booking/constant"

	"github.com/IBM/sarama"
)

type Kafka struct {
	Producer sarama.SyncProducer
}

// NewKafka tạo trước topic + khởi tạo producer dùng chung.
// Consumer group do consumer.Runner tự tạo (mỗi consumer group một instance).
func NewKafka(cfg config.KafkaConfig) (*Kafka, error) {
	if err := EnsureTopics(cfg, constant.TOPIC_PAYMENT_REQUEST); err != nil {
		return nil, err
	}

	producer, err := NewProducer(cfg)
	if err != nil {
		return nil, err
	}

	return &Kafka{Producer: producer}, nil
}

func (k *Kafka) Close() error {
	return k.Producer.Close()
}

func baseConfig() *sarama.Config {
	c := sarama.NewConfig()
	c.Version = sarama.V2_8_0_0
	return c
}

func NewProducer(cfg config.KafkaConfig) (sarama.SyncProducer, error) {
	c := baseConfig()
	c.Producer.RequiredAcks = sarama.WaitForAll
	c.Producer.Retry.Max = 3
	c.Producer.Return.Successes = true // bắt buộc với SyncProducer
	return sarama.NewSyncProducer(cfg.Brokers, c)
}

func NewConsumerGroup(cfg config.KafkaConfig, groupID string) (sarama.ConsumerGroup, error) {
	c := baseConfig()
	c.Consumer.Offsets.Initial = sarama.OffsetNewest
	c.Consumer.Return.Errors = true
	return sarama.NewConsumerGroup(cfg.Brokers, groupID, c)
}

// EnsureTopics tạo trước các topic (1 partition). Bỏ qua nếu đã tồn tại.
func EnsureTopics(cfg config.KafkaConfig, topics ...string) error {
	clusterAdmin, err := sarama.NewClusterAdmin(cfg.Brokers, baseConfig())
	if err != nil {
		return err
	}
	defer clusterAdmin.Close()

	for _, t := range topics {
		err := clusterAdmin.CreateTopic(t, &sarama.TopicDetail{
			NumPartitions:     1,
			ReplicationFactor: 1,
		}, false)
		if err != nil && !isTopicExists(err) {
			return err
		}
	}
	return nil
}

func isTopicExists(err error) bool {
	var topicErr *sarama.TopicError
	if errors.As(err, &topicErr) {
		return topicErr.Err == sarama.ErrTopicAlreadyExists
	}
	return errors.Is(err, sarama.ErrTopicAlreadyExists)
}
