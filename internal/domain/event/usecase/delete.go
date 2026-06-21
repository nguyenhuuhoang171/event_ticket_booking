package usecase

import (
	"context"
	"net/http"

	"event_ticket_booking/constant"
	eventEntity "event_ticket_booking/infrastructure/db/event/entity"
	commonModel "event_ticket_booking/model"
)

/*
1. Check event exists
2. Soft delete event
*/
func (u Usecase) Delete(ctx context.Context, userId, id uint64) error {
	// 1. Check event exists
	event, err := u.eventRepo.GetOne(ctx, eventEntity.Filter{
		Id:     id,
		Status: constant.EVENT_STATUS_ACTIVE,
	})
	if err != nil {
		return commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}
	if event == nil {
		return commonModel.NewError(http.StatusNotFound, "Event not found")
	}

	// 2. Soft delete event
	if err := u.eventRepo.Delete(ctx, &eventEntity.Entity{Id: id, DeletedBy: userId}); err != nil {
		return commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}

	return nil
}
