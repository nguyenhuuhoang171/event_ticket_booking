package usecase

import (
	"event_ticket_booking/config"
	bookingRepo "event_ticket_booking/infrastructure/db/booking/repository"
	eventRepo "event_ticket_booking/infrastructure/db/event/repository"
	commonModel "event_ticket_booking/model"

	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
)

type Usecase struct {
	cfg         config.Config
	bookingRepo bookingRepo.IRepository
	eventRepo   eventRepo.IRepository
	redis       *redis.Client
	kafka       sarama.SyncProducer
}

func NewUsecase(cfg config.Config, lib commonModel.Lib) Usecase {
	return Usecase{
		cfg:         cfg,
		bookingRepo: lib.Db.BookingRepo,
		eventRepo:   lib.Db.EventRepo,
		redis:       lib.Redis,
		kafka:       lib.KafkaProducer,
	}
}
