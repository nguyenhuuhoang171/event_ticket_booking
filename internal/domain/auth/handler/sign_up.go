package handler

import (
	"net/http"

	"event_ticket_booking/internal/domain/auth/dto"
	commonModel "event_ticket_booking/model"
	"event_ticket_booking/util"

	"github.com/gin-gonic/gin"
)

// Signup handles POST /signup.
func (h *Handler) Signup(c *gin.Context) {
	// parse request
	var request dto.SignupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.WriteError(c, commonModel.NewError(http.StatusBadRequest, util.ValidationMessage(err)))
		return
	}

	// call usecase
	signupRes, err := h.usecase.SignUp(c.Request.Context(), request)
	if err != nil {
		util.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, commonModel.Response{Data: signupRes})
}
