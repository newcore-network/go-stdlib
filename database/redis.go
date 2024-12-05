package database

import (
	"context"
	"strconv"

	"github.com/newcore-network/libs/configuration"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, cfg configuration.GeneralConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + strconv.Itoa(cfg.RedisPort),
		Password: cfg.Password,
		DB:       cfg.RedisDB,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic("cannot connect to redis")
	}

	return client
}
