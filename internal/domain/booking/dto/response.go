package dto

import (
	"time"

	bookingEntity "event_ticket_booking/infrastructure/db/booking/entity"
)

type BookingResponse struct {
	Id        uint64    `json:"id"`
	EventId   uint64    `json:"event_id"`
	UserId    uint64    `json:"user_id"`
	Quantity  uint64    `json:"quantity"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func NewBookingResponse(e *bookingEntity.Entity) BookingResponse {
	return BookingResponse{
		Id:        e.Id,
		EventId:   e.EventId,
		UserId:    e.UserId,
		Quantity:  e.Quantity,
		Status:    e.Status,
		CreatedAt: e.CreatedAt,
	}
}
