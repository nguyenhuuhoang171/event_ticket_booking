package usecase

import (
	"context"
	"log"
	"net/http"

	"event_ticket_booking/constant"
	userEntity "event_ticket_booking/infrastructure/db/user/entity"
	userRepo "event_ticket_booking/infrastructure/db/user/repository"
	"event_ticket_booking/internal/domain/auth/dto"
	commonModel "event_ticket_booking/model"
	"event_ticket_booking/util"
)

/*
1. Validate email + password
2. Check if email already exists
3. Hash password
4. Create user
5. Create access token
*/
func (u Usecase) SignUp(ctx context.Context, request dto.SignupRequest) (*dto.SignupResponse, error) {
	prefixLog := util.GetFunctionName(0)

	// 1. Validate email + password
	if !isValidEmail(request.Email) {
		return nil, commonModel.NewError(http.StatusBadRequest, "Email is invalid")
	}
	if !isValidPassword(request.Password) {
		return nil, commonModel.NewError(http.StatusBadRequest, "Password s invalid")
	}

	// 2. Check if email already exists
	user, err := u.userRepo.GetOne(ctx, userRepo.Filter{
		Email: request.Email,
	})
	if err != nil {
		log.Printf("%s Getting user: %v", prefixLog, err)
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}
	if user != nil {
		return nil, commonModel.NewError(http.StatusConflict, "Email already exists")
	}

	// 3. Hash password
	hashedPassword, err := hashPassword(request.Password)
	if err != nil {
		log.Printf("%s Hashing password: %v", prefixLog, err)
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}

	// 4. Create user
	createdUser, err := u.userRepo.Create(ctx, &userEntity.Entity{
		Email:          request.Email,
		HashedPassword: hashedPassword,
		Role:           constant.ROLE_USER,
	})
	if err != nil {
		log.Printf("%s Creating user: %v", prefixLog, err)
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}

	// 5. Create access token
	accessToken, err := u.generateAccessToken(createdUser)
	if err != nil {
		log.Printf("%s Generating access token: %v", prefixLog, err)
		return nil, commonModel.NewError(http.StatusInternalServerError, constant.INTERNAL_SERVER_ERROR)
	}

	return &dto.SignupResponse{
		AccessToken: accessToken,
	}, nil
}
