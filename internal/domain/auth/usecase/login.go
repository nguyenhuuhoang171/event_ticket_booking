package usecase

import (
	"context"
	"log"
	"net/http"

	"event_ticket_booking/constant"
	userEntity "event_ticket_booking/infrastructure/db/user/entity"
	"event_ticket_booking/internal/domain/auth/dto"
	commonModel "event_ticket_booking/model"
	"event_ticket_booking/util"
)

/*
1. Authenticate user + password
2. Create access token
*/
func (u Usecase) Login(ctx context.Context, request dto.LoginRequest) (*dto.LoginResponse, error) {
	prefixLog := util.GetFunctionName(0)

	// 1. Authenticate user + password
	user, err := u.userRepo.GetOne(ctx, userEntity.Filter{
		Email: request.Email,
	})
	if err != nil {
		log.Printf("%s Getting user: %v", prefixLog, err)
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}
	if user == nil {
		return nil, commonModel.NewError(http.StatusBadRequest, "Email not registered")
	}
	if !verifyPassword(request.Password, user.HashedPassword) {
		return nil, commonModel.NewError(http.StatusUnauthorized, "Password incorrect")
	}

	// 2. Create access token
	accessToken, err := u.generateAccessToken(user)
	if err != nil {
		log.Printf("%s Generating access token: %v", prefixLog, err)
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}

	return &dto.LoginResponse{
		AccessToken: accessToken,
	}, nil
}
