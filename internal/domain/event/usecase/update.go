package usecase

import (
	"context"
	"net/http"
	"time"

	"event_ticket_booking/constant"
	eventEntity "event_ticket_booking/infrastructure/db/event/entity"
	eventRepo "event_ticket_booking/infrastructure/db/event/repository"
	"event_ticket_booking/internal/domain/event/dto"
	commonModel "event_ticket_booking/model"
)

/*
1. Check the event exists
2. Validate params
3. Update
*/
func (u Usecase) Update(ctx context.Context, userId, id uint64, request dto.UpdateEventRequest) (*dto.UpdateResponse, error) {
	// 1. Check the event exists
	event, err := u.eventRepo.GetOne(ctx, eventRepo.Filter{
		Id:     id,
		Status: constant.EVENT_STATUS_ACTIVE,
	})
	if err != nil {
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}
	if event == nil {
		return nil, commonModel.NewError(http.StatusNotFound, "Event not found")
	}

	// 2. Validate params
	updateEvent := &eventEntity.Entity{
		Id:           id,
		UpdatedBy:    userId,
		Name:         request.Name,
		Description:  request.Description,
		TotalTickets: request.TotalTickets,
		TicketPrice:  request.TicketPrice,
	}

	// Check trùng cặp name - dateTime
	newName := event.Name
	nameChanged := request.Name != "" && request.Name != event.Name
	if nameChanged {
		newName = request.Name
	}

	newDateTime := event.DateTime
	dateTimeChanged := false
	if request.DateTime != "" && request.DateTime != event.DateTime.Format(constant.TIME_LAYOUT_YYYY_MM_DD_HH_MM_SS) {
		dateTime, err := time.ParseInLocation(constant.TIME_LAYOUT_YYYY_MM_DD_HH_MM_SS, request.DateTime, time.Local)
		if err != nil {
			return nil, commonModel.NewError(http.StatusBadRequest, "Date time phải có dạng YYYY_MM_DD hh_mm_ss")
		}
		if !dateTime.After(time.Now()) {
			return nil, commonModel.NewError(http.StatusBadRequest, "Date time phải là thời điểm trong tương lai")
		}
		newDateTime = dateTime
		updateEvent.DateTime = dateTime
		dateTimeChanged = true
	}

	if nameChanged || dateTimeChanged {
		newDateTimeStr := newDateTime.Format(constant.TIME_LAYOUT_YYYY_MM_DD_HH_MM_SS)
		if err := u.checkDuplicateEvent(ctx, id, newName, newDateTimeStr); err != nil {
			return nil, err
		}
	}

	// 3. Update
	if _, err := u.eventRepo.Update(ctx, updateEvent); err != nil {
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}
	return &dto.UpdateResponse{
		IsSuccess: true,
	}, nil
}
