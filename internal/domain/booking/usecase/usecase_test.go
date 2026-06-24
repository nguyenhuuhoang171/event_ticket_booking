package usecase

import (
	"context"
	"testing"

	"event_ticket_booking/config"
	"event_ticket_booking/constant"
	bookingEntity "event_ticket_booking/infrastructure/db/booking/entity"
	eventEntity "event_ticket_booking/infrastructure/db/event/entity"
	"event_ticket_booking/internal/domain/booking/dto"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newTestUsecase(t *testing.T, bRepo *mockBookingRepo, eRepo *mockEventRepo, kafka *mockKafkaProducer) (Usecase, *miniredis.Miniredis) {
	mr := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	u := Usecase{
		cfg:         config.Config{},
		bookingRepo: bRepo,
		eventRepo:   eRepo,
		redis:       redisClient,
		kafka:       kafka,
	}
	return u, mr
}

// Test_Create_Success: booking thành công khi event tồn tại và còn đủ vé.
func Test_Create_Success(t *testing.T) {
	bRepo := &mockBookingRepo{}
	eRepo := &mockEventRepo{}
	kafka := &mockKafkaProducer{}

	event := &eventEntity.Entity{
		Id:           1,
		TotalTickets: 100,
		SoldTickets:  10,
		TicketPrice:  50000,
		Status:       constant.EVENT_STATUS_ACTIVE,
	}
	created := &bookingEntity.Entity{
		Id:       99,
		EventId:  1,
		UserId:   7,
		Quantity: 2,
		Status:   constant.BOOKING_STATUS_PENDING,
	}

	eRepo.On("GetOne", mock.Anything, mock.Anything).Return(event, nil)
	bRepo.On("Reserve", mock.Anything, mock.Anything).Return(created, nil)
	kafka.On("SendMessage", mock.Anything).Return(int32(0), int64(0), nil)

	u, _ := newTestUsecase(t, bRepo, eRepo, kafka)
	res, err := u.Create(context.Background(), 7, dto.CreateBookingRequest{EventId: 1, Quantity: 2})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, uint64(99), res.Id)
	assert.Equal(t, constant.BOOKING_STATUS_PENDING, res.Status)
	bRepo.AssertExpectations(t)
	eRepo.AssertExpectations(t)
}

// Test_Create_SoldOut: Redis báo không đủ vé.
func Test_Create_SoldOut(t *testing.T) {
	bRepo := &mockBookingRepo{}
	eRepo := &mockEventRepo{}
	kafka := &mockKafkaProducer{}

	event := &eventEntity.Entity{
		Id:           1,
		TotalTickets: 10,
		SoldTickets:  10,
		Status:       constant.EVENT_STATUS_ACTIVE,
	}
	eRepo.On("GetOne", mock.Anything, mock.Anything).Return(event, nil)

	u, mr := newTestUsecase(t, bRepo, eRepo, kafka)
	mr.Set("event_remaining:1", "0")

	res, err := u.Create(context.Background(), 7, dto.CreateBookingRequest{EventId: 1, Quantity: 2})

	assert.Nil(t, res)
	assert.Error(t, err)
	bRepo.AssertNotCalled(t, "Reserve")
}

// Test_Cancel_Success: huỷ thành công và trả vé về Redis counter.
func Test_Cancel_Success(t *testing.T) {
	bRepo := &mockBookingRepo{}
	eRepo := &mockEventRepo{}
	kafka := &mockKafkaProducer{}

	cancelled := &bookingEntity.Entity{
		Id:       5,
		EventId:  1,
		UserId:   7,
		Quantity: 2,
		Status:   constant.BOOKING_STATUS_CANCELLED,
	}
	bRepo.On("Cancel", mock.Anything, uint64(5), uint64(7)).Return(cancelled, nil)

	u, mr := newTestUsecase(t, bRepo, eRepo, kafka)
	mr.Set("event_remaining:1", "50")

	res, err := u.Cancel(context.Background(), 7, 5)

	assert.NoError(t, err)
	assert.True(t, res.IsSuccess)

	val, _ := mr.Get("event_remaining:1")
	assert.Equal(t, "52", val)
}
