package handler

import (
	"net/http"

	"event_ticket_booking/internal/domain/auth/dto"
	commonModel "event_ticket_booking/model"
	"event_ticket_booking/util"

	"github.com/gin-gonic/gin"
)

// Logout handles POST /logout.
func (h *Handler) Logout(c *gin.Context) {
	// parse request
	var request dto.LogoutRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.WriteError(c, commonModel.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	// call usecase
	logoutRes, err := h.usecase.Logout(c.Request.Context(), request)
	if err != nil {
		util.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, commonModel.Response{Data: logoutRes})
}
