package redis

import (
	"context"
	"fmt"
	"log"

	"bj/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Redis.Password,
		DB:       0, // Use default DB
	})

	// Test connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis at %s: %v", addr, err)
		return nil, err
	}

	log.Printf("Successfully connected to Redis at %s", addr)
	return client, nil
}
