package database

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, addr string, password string, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("cannot connect to redis")
	}

	return client
}
