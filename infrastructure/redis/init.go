package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"event_ticket_booking/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisConnection(config config.RedisConfig) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Username: config.Username,
		Password: config.Password,
		DB:       config.Db,
	})

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		panic(fmt.Sprintf("failed to connect to Redis: %v", err))
	}

	log.Printf("Connect Redis successfully")
	return redisClient
}
