package util

import (
	"event_ticket_booking/constant"

	"github.com/gin-gonic/gin"
)

func GetUserId(c *gin.Context) uint64 {
	if v, ok := c.Get(constant.CONTEXT_USER_ID); ok {
		if id, ok := v.(uint64); ok {
			return id
		}
	}
	return 0
}
