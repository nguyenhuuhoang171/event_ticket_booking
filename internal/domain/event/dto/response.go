package dto

import (
	"time"

	eventEntity "event_ticket_booking/infrastructure/db/event/entity"
)

type EventResponse struct {
	Id           uint64    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	DateTime     time.Time `json:"date_time"`
	TotalTickets uint64    `json:"total_tickets"`
	TicketPrice  uint64    `json:"ticket_price"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    uint64    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    uint64    `json:"updated_by"`
}

type ListEventResponse struct {
	Items []EventResponse `json:"items"`
	Total int64           `json:"total"`
}

func NewEventResponse(e *eventEntity.Entity) EventResponse {
	return EventResponse{
		Id:           e.Id,
		Name:         e.Name,
		Description:  e.Description,
		DateTime:     e.DateTime,
		TotalTickets: e.TotalTickets,
		TicketPrice:  e.TicketPrice,
		CreatedAt:    e.CreatedAt,
		CreatedBy:    e.CreatedBy,
		UpdatedAt:    e.UpdatedAt,
		UpdatedBy:    e.UpdatedBy,
	}
}

type UpdateResponse struct {
	IsSuccess bool
}

type StatsResponse struct {
	EventId          uint64 `json:"event_id"`
	TicketsSold      uint64 `json:"tickets_sold"`
	EstimatedRevenue uint64 `json:"estimated_revenue"`
}
