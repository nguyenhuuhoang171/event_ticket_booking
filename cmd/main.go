package main

import (
	"event_ticket_booking/config"
	"event_ticket_booking/infrastructure/db"
	"event_ticket_booking/infrastructure/redis"
	commonModel "event_ticket_booking/model"
	"event_ticket_booking/server"
)

func main() {
	// init config
	cfg := config.NewConfig()

	// init lib
	lib := commonModel.Lib{
		Db:    db.NewDb(cfg),
		Redis: redis.NewRedisConnection(cfg.Redis),
	}

	// run server
	srv := server.NewServer(cfg, lib)
	defer srv.Stop()
	srv.Start()
}
