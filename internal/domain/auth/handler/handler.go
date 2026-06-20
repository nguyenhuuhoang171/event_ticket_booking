package handler

import (
	"event_ticket_booking/config"
	"event_ticket_booking/internal/domain/auth/usecase"
	commonModel "event_ticket_booking/model"
)

type Handler struct {
	cfg     config.Config
	usecase usecase.Usecase
}

func NewHandler(cfg config.Config, lib commonModel.Lib) Handler {
	authUsecase := usecase.NewUsecase(cfg, lib)
	return Handler{
		cfg:     cfg,
		usecase: authUsecase,
	}
}
