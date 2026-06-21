package usecase

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"event_ticket_booking/constant"
	"event_ticket_booking/internal/domain/auth/dto"
	"event_ticket_booking/internal/domain/auth/model"
	commonModel "event_ticket_booking/model"
	"event_ticket_booking/util"

	"github.com/golang-jwt/jwt/v5"
)

/*
1. Parse + validate access token
2. Save access token to blacklist
*/
func (u Usecase) Logout(ctx context.Context, request dto.LogoutRequest) (*dto.LogoutResponse, error) {
	prefixLog := util.GetFunctionName(0)

	// 1. Validate access token
	accessSecret := os.Getenv("ACCESS_SECRET")
	_, err := jwt.ParseWithClaims(request.AccessToken, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessSecret), nil
	})
	if err != nil {
		log.Printf("%s Parsing access token: %v", prefixLog, err)
		return nil, commonModel.NewError(http.StatusUnauthorized, "Access token is invalid")
	}

	// 2. Save access token to blacklist
	keyRedis := util.GetKeyRedis(constant.REDIS_KEY_ACCESS_TOKEN_BLACKLIST, request.AccessToken)
	if err := u.redis.Set(ctx, keyRedis, true, time.Duration(u.cfg.Authentication.AccessTokenExpirationMinutes)*time.Minute).Err(); err != nil {
		log.Printf("%s Saving access token to blacklist: %v", prefixLog, err)
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}

	return &dto.LogoutResponse{}, nil
}
