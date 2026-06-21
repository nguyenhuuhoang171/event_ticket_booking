package usecase

import (
	"context"
	"net/http"
	"time"

	"event_ticket_booking/constant"
	eventEntity "event_ticket_booking/infrastructure/db/event/entity"
	"event_ticket_booking/internal/domain/event/dto"
	commonModel "event_ticket_booking/model"
)

/*
1. Validate params
2. Check exist event with Name and Datetime if existed
3. Create new event
*/
func (u Usecase) Create(ctx context.Context, userId uint64, request dto.CreateEventRequest) (*dto.EventResponse, error) {

	// 1. Validate params
	dateTime, err := time.ParseInLocation(constant.TIME_LAYOUT_YYYY_MM_DD_HH_MM_SS, request.DateTime, time.Local)
	if err != nil {
		return nil, commonModel.NewError(http.StatusBadRequest, "Date time phải có dạng YYYY_MM_DD hh_mm_ss")
	}
	if !dateTime.After(time.Now()) {
		return nil, commonModel.NewError(http.StatusBadRequest, "Date time phải là thời điểm trong tương lai")
	}

	// 2. Check exist event with Name and Datetime if existed
	if err := u.checkDuplicateEvent(ctx, 0, request.Name, request.DateTime); err != nil {
		return nil, err
	}

	// 3. Create new event
	created, err := u.eventRepo.Create(ctx, &eventEntity.Entity{
		Name:         request.Name,
		Description:  request.Description,
		DateTime:     dateTime,
		TotalTickets: request.TotalTickets,
		TicketPrice:  request.TicketPrice,
		Status:       constant.EVENT_STATUS_ACTIVE,
		CreatedBy:    userId,
		UpdatedBy:    userId,
	})
	if err != nil {
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}

	res := dto.NewEventResponse(created)
	return &res, nil
}
