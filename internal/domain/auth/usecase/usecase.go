package usecase

import (
	"event_ticket_booking/config"
	refreshTokenRepo "event_ticket_booking/infrastructure/db/refresh_token/repository"
	userRepo "event_ticket_booking/infrastructure/db/user/repository"
	commonModel "event_ticket_booking/model"

	"github.com/redis/go-redis/v9"
)

type Usecase struct {
	cfg              config.Config
	userRepo         userRepo.IRepository
	refreshTokenRepo refreshTokenRepo.IRepository
	redis            *redis.Client
}

func NewUsecase(cfg config.Config, lib commonModel.Lib) Usecase {
	return Usecase{
		cfg:              cfg,
		userRepo:         lib.Db.UserRepo,
		refreshTokenRepo: lib.Db.RefreshTokenRepo,
		redis:            lib.Redis,
	}
}
