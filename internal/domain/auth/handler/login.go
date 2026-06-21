package handler

import (
	"net/http"

	"event_ticket_booking/internal/domain/auth/dto"
	commonModel "event_ticket_booking/model"
	"event_ticket_booking/util"

	"github.com/gin-gonic/gin"
)

// Login handles POST /login.
func (h *Handler) Login(c *gin.Context) {
	// parse request
	var request dto.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.WriteError(c, commonModel.NewError(http.StatusBadRequest, util.ValidationMessage(err)))
		return
	}

	// call usecase
	loginRes, err := h.usecase.Login(c.Request.Context(), request)
	if err != nil {
		util.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, commonModel.Response{Data: loginRes})
}
