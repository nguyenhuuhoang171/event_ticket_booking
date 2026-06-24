package usecase

import (
	"context"

	"event_ticket_booking/infrastructure/db/base"
	bookingEntity "event_ticket_booking/infrastructure/db/booking/entity"
	bookingRepo "event_ticket_booking/infrastructure/db/booking/repository"
	eventEntity "event_ticket_booking/infrastructure/db/event/entity"
	eventRepo "event_ticket_booking/infrastructure/db/event/repository"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/mock"
)

// booking repo mock

type mockBookingRepo struct{ mock.Mock }

func (m *mockBookingRepo) Reserve(ctx context.Context, e *bookingEntity.Entity) (*bookingEntity.Entity, error) {
	args := m.Called(ctx, e)
	v, _ := args.Get(0).(*bookingEntity.Entity)
	return v, args.Error(1)
}
func (m *mockBookingRepo) Cancel(ctx context.Context, bookingId, userId uint64) (*bookingEntity.Entity, error) {
	args := m.Called(ctx, bookingId, userId)
	v, _ := args.Get(0).(*bookingEntity.Entity)
	return v, args.Error(1)
}
func (m *mockBookingRepo) GetListPaging(ctx context.Context, filter bookingRepo.Filter, page, size int64) ([]bookingEntity.Entity, int64, error) {
	return nil, 0, nil
}

func (m *mockBookingRepo) Create(ctx context.Context, e *bookingEntity.Entity) (*bookingEntity.Entity, error) {
	return nil, nil
}
func (m *mockBookingRepo) CreateMany(ctx context.Context, entities []*bookingEntity.Entity) ([]*bookingEntity.Entity, error) {
	return nil, nil
}
func (m *mockBookingRepo) GetOne(ctx context.Context, filter bookingRepo.Filter, opts ...base.QueryOption) (*bookingEntity.Entity, error) {
	return nil, nil
}
func (m *mockBookingRepo) GetMany(ctx context.Context, filter bookingRepo.Filter, opts ...base.QueryOption) ([]bookingEntity.Entity, error) {
	return nil, nil
}
func (m *mockBookingRepo) Count(ctx context.Context, filter bookingRepo.Filter) (int64, error) {
	return 0, nil
}
func (m *mockBookingRepo) Update(ctx context.Context, e *bookingEntity.Entity) (*bookingEntity.Entity, error) {
	return nil, nil
}
func (m *mockBookingRepo) UpdateMany(ctx context.Context, entities []*bookingEntity.Entity) ([]*bookingEntity.Entity, error) {
	return nil, nil
}
func (m *mockBookingRepo) Delete(ctx context.Context, e *bookingEntity.Entity) error { return nil }
func (m *mockBookingRepo) Confirm(ctx context.Context, bookingId uint64) (*bookingEntity.Entity, error) {
	return nil, nil
}
func (m *mockBookingRepo) CancelBookings(ctx context.Context, bookingIds []uint64) ([]bookingEntity.Entity, error) {
	return nil, nil
}
func (m *mockBookingRepo) GetStats(ctx context.Context, eventId uint64) (*bookingEntity.Entity, error) {
	return nil, nil
}

// event repo mock

type mockEventRepo struct{ mock.Mock }

func (m *mockEventRepo) GetOne(ctx context.Context, filter eventRepo.Filter, opts ...base.QueryOption) (*eventEntity.Entity, error) {
	args := m.Called(ctx, filter)
	v, _ := args.Get(0).(*eventEntity.Entity)
	return v, args.Error(1)
}

func (m *mockEventRepo) Create(ctx context.Context, e *eventEntity.Entity) (*eventEntity.Entity, error) {
	return nil, nil
}
func (m *mockEventRepo) CreateMany(ctx context.Context, entities []*eventEntity.Entity) ([]*eventEntity.Entity, error) {
	return nil, nil
}
func (m *mockEventRepo) GetMany(ctx context.Context, filter eventRepo.Filter, opts ...base.QueryOption) ([]eventEntity.Entity, error) {
	return nil, nil
}
func (m *mockEventRepo) Count(ctx context.Context, filter eventRepo.Filter) (int64, error) {
	return 0, nil
}
func (m *mockEventRepo) Update(ctx context.Context, e *eventEntity.Entity) (*eventEntity.Entity, error) {
	return nil, nil
}
func (m *mockEventRepo) UpdateMany(ctx context.Context, entities []*eventEntity.Entity) ([]*eventEntity.Entity, error) {
	return nil, nil
}
func (m *mockEventRepo) Delete(ctx context.Context, e *eventEntity.Entity) error { return nil }
func (m *mockEventRepo) GetListPaging(ctx context.Context, filter eventRepo.Filter, page, size int) ([]eventEntity.Entity, int64, error) {
	return nil, 0, nil
}

// kafka producer mock

type mockKafkaProducer struct{ mock.Mock }

func (m *mockKafkaProducer) SendMessage(msg *sarama.ProducerMessage) (int32, int64, error) {
	args := m.Called(msg)
	return args.Get(0).(int32), args.Get(1).(int64), args.Error(2)
}

func (m *mockKafkaProducer) SendMessages(msgs []*sarama.ProducerMessage) error { return nil }
func (m *mockKafkaProducer) Close() error                                      { return nil }
func (m *mockKafkaProducer) TxnStatus() sarama.ProducerTxnStatusFlag           { return 0 }
func (m *mockKafkaProducer) IsTransactional() bool                             { return false }
func (m *mockKafkaProducer) BeginTxn() error                                   { return nil }
func (m *mockKafkaProducer) CommitTxn() error                                  { return nil }
func (m *mockKafkaProducer) AbortTxn() error                                   { return nil }
func (m *mockKafkaProducer) AddMessageToTxn(_ *sarama.ConsumerMessage, _ string, _ *string) error {
	return nil
}
func (m *mockKafkaProducer) AddMessageToTxnWithGroupMetadata(_ *sarama.ConsumerMessage, _ *sarama.ConsumerGroupMetadata, _ *string) error {
	return nil
}
func (m *mockKafkaProducer) AddOffsetsToTxn(_ map[string][]*sarama.PartitionOffsetMetadata, _ string) error {
	return nil
}
func (m *mockKafkaProducer) AddOffsetsToTxnWithGroupMetadata(_ map[string][]*sarama.PartitionOffsetMetadata, _ *sarama.ConsumerGroupMetadata) error {
	return nil
}
