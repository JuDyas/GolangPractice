package database

import (
	"github.com/go-redis/redis/v8"
)

func SetupRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
