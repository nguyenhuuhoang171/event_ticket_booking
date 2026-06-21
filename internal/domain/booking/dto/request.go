package dto

type CreateBookingRequest struct {
	EventId  uint64 `json:"event_id" binding:"required"`
	Quantity uint64 `json:"quantity" binding:"required,min=1"`
}
