package dto

type CreateEventRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	DateTime     string `json:"date_time" binding:"required"`
	TotalTickets uint64 `json:"total_tickets" binding:"required"`
	TicketPrice  uint64 `json:"ticket_price" binding:"required"`
}

type UpdateEventRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	DateTime     string `json:"date_time"`
	TotalTickets uint64 `json:"total_tickets"`
	TicketPrice  uint64 `json:"ticket_price"`
}

type ListEventRequest struct {
	Name string `form:"name"`
	Page int    `form:"page"`
	Size int    `form:"size"`
}

type StatsRequest struct {
	EventId uint64 `form:"event_id"`
}
