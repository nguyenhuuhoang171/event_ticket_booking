package usecase

import (
	"context"
	"net/http"

	"event_ticket_booking/constant"
	bookingRepo "event_ticket_booking/infrastructure/db/booking/repository"
	"event_ticket_booking/internal/domain/booking/dto"
	commonModel "event_ticket_booking/model"
)

func (u Usecase) List(ctx context.Context, userId uint64, request dto.ListBookingRequest) (*dto.ListBookingResponse, error) {
	filter := bookingRepo.Filter{
		UserId:  userId,
		EventId: request.EventId,
		Status:  request.Status,
	}
	items, total, err := u.bookingRepo.GetList(ctx, filter, request.Page, request.Size)
	if err != nil {
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}

	res := dto.ListBookingResponse{
		Items: make([]dto.BookingResponse, 0, len(items)),
		Total: total,
	}
	for i := range items {
		res.Items = append(res.Items, dto.NewBookingResponse(&items[i]))
	}

	return &res, nil
}
