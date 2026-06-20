package usecase

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"event_ticket_booking/constant"
	refreshTokenEntity "event_ticket_booking/infrastructure/db/refresh_token/entity"
	"event_ticket_booking/internal/domain/auth/dto"
	"event_ticket_booking/internal/domain/auth/model"
	commonModel "event_ticket_booking/model"
	"event_ticket_booking/util"

	"github.com/golang-jwt/jwt/v5"
)

/*
1. Parse access token to identify the user
2. Save access token to blacklist
3. Revoke the user's refresh tokens
*/
func (u Usecase) Logout(ctx context.Context, request dto.LogoutRequest) (*dto.LogoutResponse, error) {
	prefixLog := util.GetFunctionName(0)

	// 1. Parse access token to identify the user
	accessSecret := os.Getenv("ACCESS_SECRET")
	token, err := jwt.ParseWithClaims(request.AccessToken, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessSecret), nil
	})
	if err != nil {
		log.Printf("%s Parsing access token: %v", prefixLog, err)
		return nil, commonModel.NewError(http.StatusUnauthorized, "Access token is invalid")
	}
	claims, ok := token.Claims.(*model.Claims)
	if !ok {
		log.Printf("%s Token claims are invalid", prefixLog)
		return nil, commonModel.NewError(http.StatusUnauthorized, "Access token is invalid")
	}

	// 2. Save access token to blacklist
	keyRedis := util.GetKeyRedis(constant.REDIS_KEY_ACCESS_TOKEN_BLACKLIST, request.AccessToken)
	err = u.redis.Set(ctx, keyRedis, true, time.Duration(u.cfg.Authentication.AccessTokenExpirationMinutes)*time.Minute).Err()
	if err != nil {
		log.Printf("%s Saving access token to blacklist: %v", prefixLog, err)
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}

	// 3. Revoke the user's refresh tokens
	err = u.refreshTokenRepo.Revoke(ctx, refreshTokenEntity.Filter{
		UserId: claims.UserId,
	})
	if err != nil {
		log.Printf("%s Revoking refresh token: %v", prefixLog, err)
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}

	return &dto.LogoutResponse{}, nil
}
