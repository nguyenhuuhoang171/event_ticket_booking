package usecase

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"event_ticket_booking/constant"
	bookingEntity "event_ticket_booking/infrastructure/db/booking/entity"
	bookingRepo "event_ticket_booking/infrastructure/db/booking/repository"
	eventRepo "event_ticket_booking/infrastructure/db/event/repository"
	"event_ticket_booking/internal/domain/booking/dto"
	commonModel "event_ticket_booking/model"
	"event_ticket_booking/util"

	"github.com/redis/go-redis/v9"
)

// KEYS[1] = remaining counter key
// ARGV[1] = quantity
// ARGV[2] = number of available tickets
var reserveScript = redis.NewScript(`
if redis.call('EXISTS', KEYS[1]) == 0 then
	redis.call('SET', KEYS[1], ARGV[2])
end
local remaining = tonumber(redis.call('GET', KEYS[1]))
local qty = tonumber(ARGV[1])
if remaining < qty then
	return -1
end
return redis.call('DECRBY', KEYS[1], qty)
`)

/*
1. Check event exists
2. Reserve tickets atomically on Redis
3. Create booking
4. Roll back the Redis if transaction has error
*/
func (u Usecase) Create(ctx context.Context, userId uint64, request dto.CreateBookingRequest) (*dto.BookingResponse, error) {
	prefixLog := util.GetFunctionName(0)

	// 1. Check event exists
	event, err := u.eventRepo.GetOne(ctx, eventRepo.Filter{
		Id:     request.EventId,
		Status: constant.EVENT_STATUS_ACTIVE,
	})
	if err != nil {
		log.Printf("%s Getting event: %v", prefixLog, err)
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}
	if event == nil {
		return nil, commonModel.NewError(http.StatusNotFound, "Event not found")
	}

	// 2. Reserve tickets atomically on Redis
	remainingKey := util.GetKeyRedis(constant.REDIS_KEY_EVENT_REMAINING, strconv.FormatUint(event.Id, 10))
	available := event.TotalTickets - event.SoldTickets
	reserved, err := reserveScript.Run(ctx, u.redis, []string{remainingKey}, request.Quantity, available).Int64()
	if err != nil {
		log.Printf("%s Reserving tickets on redis: %v", prefixLog, err)
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}
	if reserved < 0 {
		return nil, commonModel.NewError(http.StatusConflict, "Không đủ vé")
	}

	// 3. Create booking
	booking := &bookingEntity.Entity{
		EventId:   request.EventId,
		UserId:    userId,
		Quantity:  request.Quantity,
		Status:    constant.BOOKING_STATUS_PENDING,
		CreatedBy: userId,
		UpdatedBy: userId,
	}
	created, err := u.bookingRepo.Reserve(ctx, booking)
	if err != nil {
		// 4. Roll back the Redis if transaction has error
		if rbErr := u.redis.IncrBy(ctx, remainingKey, int64(request.Quantity)).Err(); rbErr != nil {
			log.Printf("%s Rolling back redis reservation: %v", prefixLog, rbErr)
		}
		if errors.Is(err, bookingRepo.ErrSoldOut) {
			return nil, commonModel.NewError(http.StatusConflict, "Không đủ vé")
		}
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}

	res := dto.NewBookingResponse(created)
	return &res, nil
}
