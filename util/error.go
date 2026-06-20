package util

import (
	"errors"
	"net/http"

	commonModel "event_ticket_booking/model"

	"github.com/gin-gonic/gin"
)

func WriteError(c *gin.Context, err error) {
	var appErr *commonModel.Error
	if errors.As(err, &appErr) {
		c.JSON(appErr.Status, commonModel.Response{Error: appErr})
		return
	}

	c.JSON(http.StatusInternalServerError, commonModel.Response{
		Error: &commonModel.Error{Message: err.Error()},
	})
}
