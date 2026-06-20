package model

import (
	"event_ticket_booking/infrastructure/db"

	"github.com/redis/go-redis/v9"
)

type Lib struct {
	Db    db.Db
	Redis *redis.Client
}
