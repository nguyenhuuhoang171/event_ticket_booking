package handler

import (
	"net/http"
	"strconv"

	"event_ticket_booking/config"
	"event_ticket_booking/internal/domain/event/dto"
	"event_ticket_booking/internal/domain/event/usecase"
	commonModel "event_ticket_booking/model"
	"event_ticket_booking/util"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	cfg     config.Config
	usecase usecase.Usecase
}

func NewHandler(cfg config.Config, lib commonModel.Lib) Handler {
	eventUsecase := usecase.NewUsecase(cfg, lib)
	return Handler{
		cfg:     cfg,
		usecase: eventUsecase,
	}
}

func (h *Handler) Create(c *gin.Context) {
	var request dto.CreateEventRequest
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
	var request dto.ListEventRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		util.WriteError(c, commonModel.NewError(http.StatusBadRequest, util.ValidationMessage(err)))
		return
	}

	res, err := h.usecase.List(c.Request.Context(), request)
	if err != nil {
		util.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, commonModel.Response{Data: res})
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.WriteError(c, commonModel.NewError(http.StatusBadRequest, "Event id is invalid"))
		return
	}

	res, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		util.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, commonModel.Response{Data: res})
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.WriteError(c, commonModel.NewError(http.StatusBadRequest, "Event id is invalid"))
		return
	}

	var request dto.UpdateEventRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.WriteError(c, commonModel.NewError(http.StatusBadRequest, util.ValidationMessage(err)))
		return
	}

	res, err := h.usecase.Update(c.Request.Context(), util.GetUserId(c), id, request)
	if err != nil {
		util.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, commonModel.Response{Data: res})
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.WriteError(c, commonModel.NewError(http.StatusBadRequest, "Event id is invalid"))
		return
	}

	if err := h.usecase.Delete(c.Request.Context(), util.GetUserId(c), id); err != nil {
		util.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, commonModel.Response{})
}

func (h *Handler) Stats(c *gin.Context) {
	var eventId uint64
	if raw := c.Query("event_id"); raw != "" {
		id, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			util.WriteError(c, commonModel.NewError(http.StatusBadRequest, "Event id is invalid"))
			return
		}
		eventId = id
	}

	res, err := h.usecase.GetStats(c.Request.Context(), eventId)
	if err != nil {
		util.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, commonModel.Response{Data: res})
}
