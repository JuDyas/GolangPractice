package config

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func SetupRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
