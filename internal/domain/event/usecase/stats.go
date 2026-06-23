package usecase

import (
	"context"
	"net/http"

	"event_ticket_booking/constant"
	eventRepo "event_ticket_booking/infrastructure/db/event/repository"
	"event_ticket_booking/internal/domain/event/dto"
	commonModel "event_ticket_booking/model"
)

func (u Usecase) GetStats(ctx context.Context, eventId uint64) (*dto.StatsResponse, error) {

	event, err := u.eventRepo.GetOne(ctx, eventRepo.Filter{
		Id:     eventId,
		Status: constant.EVENT_STATUS_ACTIVE,
	})
	if err != nil {
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}
	if event == nil {
		return nil, commonModel.NewError(http.StatusNotFound, "Event not found")
	}

	stats, err := u.bookingRepo.GetStats(ctx, eventId)
	if err != nil {
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}

	return &dto.StatsResponse{
		EventId:          eventId,
		TicketsSold:      stats.TicketsSold,
		EstimatedRevenue: stats.EstimatedRevenue,
	}, nil
}
