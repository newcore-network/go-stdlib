package database

import (
	"context"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/styerr-development/libs/configuration"
)

func NewRedisClient(ctx context.Context, cfg configuration.GeneralConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + strconv.Itoa(cfg.RedisPort),
		Password: cfg.Password,
		DB:       cfg.RedisDB,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("cannot connect to redis")
	}

	return client
}
