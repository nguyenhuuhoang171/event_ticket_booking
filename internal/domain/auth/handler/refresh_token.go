package handler

import (
	"net/http"

	"event_ticket_booking/internal/domain/auth/dto"
	commonModel "event_ticket_booking/model"
	"event_ticket_booking/util"

	"github.com/gin-gonic/gin"
)

// RefreshToken handles POST /refresh-token.
func (h *Handler) RefreshToken(c *gin.Context) {
	// parse request
	var request dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.WriteError(c, commonModel.NewError(http.StatusBadRequest, err.Error()))
		return
	}

	// call usecase
	refreshTokenRes, err := h.usecase.RefreshToken(c.Request.Context(), request)
	if err != nil {
		util.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, commonModel.Response{Data: refreshTokenRes})
}
