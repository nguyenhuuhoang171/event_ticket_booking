package dto

type CreateBookingRequest struct {
	EventId  uint64 `json:"event_id" binding:"required"`
	Quantity uint64 `json:"quantity" binding:"required,min=1"`
}

type ListBookingRequest struct {
	EventId uint64 `form:"event_id"`
	Status  int    `form:"status"`
	Page    int    `form:"page"`
	Size    int    `form:"size"`
}
