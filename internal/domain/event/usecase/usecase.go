package usecase

import (
	"context"
	"event_ticket_booking/config"
	"event_ticket_booking/constant"
	bookingRepo "event_ticket_booking/infrastructure/db/booking/repository"
	eventRepo "event_ticket_booking/infrastructure/db/event/repository"
	commonModel "event_ticket_booking/model"
	"fmt"
	"net/http"
)

type Usecase struct {
	cfg         config.Config
	eventRepo   eventRepo.IRepository
	bookingRepo bookingRepo.IRepository
}

func NewUsecase(cfg config.Config, lib commonModel.Lib) Usecase {
	return Usecase{
		cfg:         cfg,
		eventRepo:   lib.Db.EventRepo,
		bookingRepo: lib.Db.BookingRepo,
	}
}

func (u Usecase) checkDuplicateEvent(ctx context.Context, excludeId uint64, name, dateTimeStr string) error {
	existedEvent, err := u.eventRepo.GetOne(ctx, eventRepo.Filter{
		Name:     name,
		DateTime: dateTimeStr,
		Status:   constant.EVENT_STATUS_ACTIVE,
	})
	if err != nil {
		return commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}
	if existedEvent != nil && existedEvent.Id != excludeId {
		return commonModel.NewError(http.StatusBadRequest, fmt.Sprintf("Sự kiện %v vào lúc %v đã được tạo", name, dateTimeStr))
	}
	return nil
}
