package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"event_ticket_booking/config"
	"event_ticket_booking/infrastructure/db"
	kafkaInfra "event_ticket_booking/infrastructure/kafka"
	"event_ticket_booking/infrastructure/redis"
	"event_ticket_booking/internal/worker"
	commonModel "event_ticket_booking/model"
	"event_ticket_booking/server"
)

func main() {
	// init config
	cfg := config.NewConfig()

	// init kafka (topic + producer)
	kafka, err := kafkaInfra.NewKafka(cfg.Kafka)
	if err != nil {
		log.Fatalf("Fail to init kafka: %v", err)
	}

	// init lib
	lib := commonModel.Lib{
		Db:            db.NewDb(cfg),
		Redis:         redis.NewRedisConnection(cfg.Redis),
		KafkaProducer: kafka.Producer,
	}

	ctx, cancel := context.WithCancel(context.Background())

	// start worker background: consumer Kafka + cron
	workers, err := worker.InitWorker(ctx, cfg, lib)
	if err != nil {
		log.Fatalf("Fail to start worker: %v", err)
	}
	log.Printf("Payment consumer & timeout cron started")

	// run HTTP server
	srv := server.NewServer(cfg, lib)
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Fail to start server: %v", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("Shutting down...")
	cancel()
	workers.Stop()
	kafka.Close()
	srv.Stop()
}
