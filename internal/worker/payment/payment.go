package payment

import (
	"context"
	"log"
	"strconv"
	"time"

	"event_ticket_booking/config"
	"event_ticket_booking/constant"
	bookingEntity "event_ticket_booking/infrastructure/db/booking/entity"
	bookingRepo "event_ticket_booking/infrastructure/db/booking/repository"
	eventRepo "event_ticket_booking/infrastructure/db/event/repository"
	commonModel "event_ticket_booking/model"
	"event_ticket_booking/util"

	"github.com/redis/go-redis/v9"
)

type Processor struct {
	cfg         config.Config
	bookingRepo bookingRepo.IRepository
	eventRepo   eventRepo.IRepository
	redis       *redis.Client
}

func NewProcessor(cfg config.Config, lib commonModel.Lib) Processor {
	return Processor{
		cfg:         cfg,
		bookingRepo: lib.Db.BookingRepo,
		eventRepo:   lib.Db.EventRepo,
		redis:       lib.Redis,
	}
}

func (p Processor) KafkaConfig() config.KafkaConfig {
	return p.cfg.Kafka
}

/*
ProcessPayment mô phỏng xử lý thanh toán cho 1 booking
- Bỏ qua nếu booking không còn PENDING.
- Mô phỏng thanh toán: thành công -> CONFIRMED, thất bại -> CANCELLED + trả vé.
*/
func (p Processor) ProcessPayment(ctx context.Context, bookingId uint64) error {
	prefixLog := util.GetFunctionName(0)

	booking, err := p.bookingRepo.GetOne(ctx, bookingRepo.Filter{Id: bookingId})
	if err != nil {
		log.Printf("%s Getting booking: %v", prefixLog, err)
		return err
	}
	if booking == nil || booking.Status != constant.BOOKING_STATUS_PENDING {
		return nil
	}

	if err := simulatePay(bookingId); err != nil {
		p.FailPayment(ctx, []bookingEntity.Entity{*booking})
		return nil
	}

	if _, err := p.bookingRepo.Update(ctx, &bookingEntity.Entity{
		Id:        bookingId,
		Status:    constant.BOOKING_STATUS_CONFIRMED,
		UpdatedBy: constant.SYSTEM_USER_ID,
	}); err != nil {
		log.Printf("%s Confirming booking: %v", prefixLog, err)
		return err
	}
	log.Printf("%s Booking %d payment SUCCESS -> confirmed", prefixLog, bookingId)
	return nil
}

// GetExpiredPendingBookings lấy 1 batch (tối đa MAX_SIZE) booking còn PENDING đã quá hạn.
func (p Processor) GetExpiredPendingBookings(ctx context.Context) ([]bookingEntity.Entity, error) {
	before := time.Now().Add(-time.Duration(p.cfg.Payment.TimeoutMinutes) * time.Minute)
	bookings, _, err := p.bookingRepo.GetListPaging(ctx, bookingRepo.Filter{
		Status:      constant.BOOKING_STATUS_PENDING,
		ToCreatedAt: before,
	}, 1, constant.MAX_SIZE)
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

// FailPayment huỷ loạt booking PENDING + trả vé về Redis counter.
func (p Processor) FailPayment(ctx context.Context, bookings []bookingEntity.Entity) error {
	prefixLog := util.GetFunctionName(0)

	if len(bookings) == 0 {
		return nil
	}

	bookingIds := make([]uint64, len(bookings))
	for i := range bookings {
		bookingIds[i] = bookings[i].Id
	}

	cancelled, err := p.bookingRepo.CancelBookings(ctx, bookingIds)
	if err != nil {
		log.Printf("%s Cancelling bookings: %v", prefixLog, err)
		return err
	}

	for i := range cancelled {
		remainingKey := util.GetKeyRedis(constant.REDIS_KEY_EVENT_REMAINING, strconv.FormatUint(cancelled[i].EventId, 10))
		if err := redis.NewScript(constant.REDIS_SCRIPT_RELEASE_TICKETS).Run(ctx, p.redis, []string{remainingKey}, int64(cancelled[i].Quantity)).Err(); err != nil {
			log.Printf("%s Releasing tickets to redis %d: %v", prefixLog, cancelled[i].Id, err)
		}
	}
	return nil
}

func simulatePay(bookingId uint64) error {
	log.Printf("[Payment] Paying booking %d...", bookingId)
	return nil
}
