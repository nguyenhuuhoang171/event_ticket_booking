package usecase

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"event_ticket_booking/constant"
	bookingRepo "event_ticket_booking/infrastructure/db/booking/repository"
	"event_ticket_booking/internal/domain/booking/dto"
	commonModel "event_ticket_booking/model"
	"event_ticket_booking/util"

	"github.com/redis/go-redis/v9"
)

// KEYS[1] = remaining counter key
// ARGV[1] = quantity to release
var releaseScript = redis.NewScript(`
if redis.call('EXISTS', KEYS[1]) == 1 then
	return redis.call('INCRBY', KEYS[1], ARGV[1])
end
return -1
`)

/*
1. Cancel booking + release tickets in DB
2. Release tickets back to the Redis counter.
*/
func (u Usecase) Cancel(ctx context.Context, userId, bookingId uint64) (*dto.CancelResponse, error) {
	prefixLog := util.GetFunctionName(0)

	// 1. Cancel booking + release tickets in DB
	cancelled, err := u.bookingRepo.Cancel(ctx, bookingId, userId)
	if err != nil {
		switch {
		case errors.Is(err, bookingRepo.ErrBookingNotFound):
			return nil, commonModel.NewError(http.StatusNotFound, "Không tìm thấy booking")
		case errors.Is(err, bookingRepo.ErrAlreadyCancelled):
			return nil, commonModel.NewError(http.StatusConflict, "Booking đã bị huỷ")
		default:
			log.Printf("%s Cancelling booking: %v", prefixLog, err)
			return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
		}
	}

	// 2. Release tickets back to the Redis counter.
	remainingKey := util.GetKeyRedis(constant.REDIS_KEY_EVENT_REMAINING, strconv.FormatUint(cancelled.EventId, 10))
	if err := releaseScript.Run(ctx, u.redis, []string{remainingKey}, int64(cancelled.Quantity)).Err(); err != nil {
		log.Printf("%s Releasing tickets to redis: %v", prefixLog, err)
	}

	return &dto.CancelResponse{IsSuccess: true}, nil
}
