package usecase

import (
	"context"
	"log"
	"net/http"

	"event_ticket_booking/constant"
	eventEntity "event_ticket_booking/infrastructure/db/event/entity"
	"event_ticket_booking/internal/domain/event/dto"
	commonModel "event_ticket_booking/model"
	"event_ticket_booking/util"
)

func (u Usecase) GetByID(ctx context.Context, id uint64) (*dto.EventResponse, error) {
	prefixLog := util.GetFunctionName(0)

	event, err := u.eventRepo.GetOne(ctx, eventEntity.Filter{
		Id:     id,
		Status: constant.EVENT_STATUS_ACTIVE,
	})
	if err != nil {
		log.Printf("%s Getting event: %v", prefixLog, err)
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}
	if event == nil {
		return nil, commonModel.NewError(http.StatusNotFound, "Event not found")
	}

	res := dto.NewEventResponse(event)
	return &res, nil
}

func (u Usecase) List(ctx context.Context, request dto.ListEventRequest) (*dto.ListEventResponse, error) {
	prefixLog := util.GetFunctionName(0)

	items, total, err := u.eventRepo.GetList(ctx, eventEntity.Filter{Name: request.Name}, request.Page, request.Size)
	if err != nil {
		log.Printf("%s Listing events: %v", prefixLog, err)
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}

	res := dto.ListEventResponse{
		Items: make([]dto.EventResponse, 0, len(items)),
		Total: total,
	}
	for i := range items {
		res.Items = append(res.Items, dto.NewEventResponse(&items[i]))
	}

	return &res, nil
}
