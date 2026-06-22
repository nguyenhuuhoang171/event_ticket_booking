package handler

import (
	"net/http"
	"strconv"

	"event_ticket_booking/config"
	"event_ticket_booking/internal/domain/booking/dto"
	"event_ticket_booking/internal/domain/booking/usecase"
	commonModel "event_ticket_booking/model"
	"event_ticket_booking/util"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	cfg     config.Config
	usecase usecase.Usecase
}

func NewHandler(cfg config.Config, lib commonModel.Lib) Handler {
	bookingUsecase := usecase.NewUsecase(cfg, lib)
	return Handler{
		cfg:     cfg,
		usecase: bookingUsecase,
	}
}

func (h *Handler) Create(c *gin.Context) {
	var request dto.CreateBookingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.WriteError(c, commonModel.NewError(http.StatusBadRequest, util.ValidationMessage(err)))
		return
	}

	res, err := h.usecase.Create(c.Request.Context(), util.GetUserId(c), request)
	if err != nil {
		util.WriteError(c, err)
		return
	}

	c.JSON(http.StatusCreated, commonModel.Response{Data: res})
}

func (h *Handler) List(c *gin.Context) {
	var request dto.ListBookingRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		util.WriteError(c, commonModel.NewError(http.StatusBadRequest, util.ValidationMessage(err)))
		return
	}

	res, err := h.usecase.List(c.Request.Context(), util.GetUserId(c), request)
	if err != nil {
		util.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, commonModel.Response{Data: res})
}

func (h *Handler) Cancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.WriteError(c, commonModel.NewError(http.StatusBadRequest, "Booking id is invalid"))
		return
	}

	res, err := h.usecase.Cancel(c.Request.Context(), util.GetUserId(c), id)
	if err != nil {
		util.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, commonModel.Response{Data: res})
}
