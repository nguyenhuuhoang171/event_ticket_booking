package ping

import (
	"event_ticket_booking/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping handles GET /ping for health checks.
func Ping(c *gin.Context) {
	response := model.Response{
		Data: "pong",
	}
	c.JSON(http.StatusOK, response)
}
