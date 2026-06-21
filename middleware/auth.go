package middleware

import (
	"net/http"
	"strings"

	"event_ticket_booking/constant"
	authModel "event_ticket_booking/internal/domain/auth/model"
	commonModel "event_ticket_booking/model"
	"event_ticket_booking/util"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func Authorize(accessSecret string, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if token == "" || token == authHeader {
			util.WriteError(c, commonModel.NewError(http.StatusUnauthorized, "Access token is required"))
			c.Abort()
			return
		}

		// Parse claim
		claims := &authModel.Claims{}
		_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(accessSecret), nil
		})
		if err != nil {
			util.WriteError(c, commonModel.NewError(http.StatusUnauthorized, "Access token is invalid"))
			c.Abort()
			return
		}

		// Reject tokens that have been logged out.
		keyRedis := util.GetKeyRedis(constant.REDIS_KEY_ACCESS_TOKEN_BLACKLIST, token)
		if exists, err := redisClient.Exists(c.Request.Context(), keyRedis).Result(); err != nil {
			util.WriteError(c, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR))
			c.Abort()
			return
		} else if exists > 0 {
			util.WriteError(c, commonModel.NewError(http.StatusUnauthorized, "Access token is revoked"))
			c.Abort()
			return
		}

		c.Set(constant.CONTEXT_USER_ID, claims.UserId)
		c.Next()
	}
}
