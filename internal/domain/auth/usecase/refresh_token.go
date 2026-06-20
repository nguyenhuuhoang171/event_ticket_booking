package usecase

import (
	"context"
	"log"
	"net/http"
	"time"

	"event_ticket_booking/constant"
	refreshTokenEntity "event_ticket_booking/infrastructure/db/refresh_token/entity"
	userEntity "event_ticket_booking/infrastructure/db/user/entity"
	"event_ticket_booking/internal/domain/auth/dto"
	commonModel "event_ticket_booking/model"
	"event_ticket_booking/util"
)

/*
1. Validate refresh token (check existence, active, not expired)
2. Revoke the old refresh token
3. Get user from refresh token
4. Generate new access token
5. Generate new refresh token (rotate)
6. Add old access token to blacklist
*/
func (u Usecase) RefreshToken(ctx context.Context, request dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	prefixLog := util.GetFunctionName(0)

	// 1. Validate refresh token (check existence, active, not expired)
	refreshTokenRecord, err := u.refreshTokenRepo.GetOne(ctx, refreshTokenEntity.Filter{
		Token:  request.RefreshToken,
		Status: constant.REFRESH_TOKEN_STATUS_ACTIVE,
	})
	if err != nil {
		log.Printf("%s Getting refresh token: %v", prefixLog, err)
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}
	if refreshTokenRecord == nil {
		return nil, commonModel.NewError(http.StatusUnauthorized, "Refresh token is invalid")
	}
	if refreshTokenRecord.ExpireAt.Before(time.Now()) {
		return nil, commonModel.NewError(http.StatusUnauthorized, "Refresh token has expired")
	}

	// 2. Revoke the old refresh token
	revokedStatus := constant.REFRESH_TOKEN_STATUS_REVOKED
	_, err = u.refreshTokenRepo.Update(ctx, &refreshTokenEntity.Entity{
		Id:     refreshTokenRecord.Id,
		Status: &revokedStatus,
	})
	if err != nil {
		log.Printf("%s Revoking refresh token: %v", prefixLog, err)
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}

	// 3. Get user from refresh token
	user, err := u.userRepo.GetOne(ctx, userEntity.Filter{
		Id: refreshTokenRecord.UserId,
	})
	if err != nil {
		log.Printf("%s Getting user: %v", prefixLog, err)
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}
	if user == nil {
		return nil, commonModel.NewError(http.StatusNotFound, "User not found")
	}

	// 4. Generate new access token
	accessToken, err := u.generateAccessToken(user)
	if err != nil {
		log.Printf("%s Generating access token: %v", prefixLog, err)
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}

	// 5. Generate new refresh token (rotate)
	newRefreshToken, err := u.saveRefreshToken(ctx, user.Id)
	if err != nil {
		log.Printf("%s Saving new refresh token: %v", prefixLog, err)
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}

	// 6. Add old access token to blacklist
	keyRedis := util.GetKeyRedis(constant.REDIS_KEY_ACCESS_TOKEN_BLACKLIST, request.AccessToken)
	err = u.redis.Set(ctx, keyRedis, true, time.Duration(u.cfg.Authentication.AccessTokenExpirationMinutes)*time.Minute).Err()
	if err != nil {
		log.Printf("%s Adding access token to blacklist: %v", prefixLog, err)
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}

	return &dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
