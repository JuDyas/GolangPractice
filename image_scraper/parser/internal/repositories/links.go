package repositories

import (
	"context"
	"fmt"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/parser/internal/models"

	"github.com/redis/go-redis/v9"
)

type LinksRepository interface {
	SaveLink(ctx context.Context, url string)
	CheckExists(ctx context.Context, url string) bool
}

type LinksRepositoryImpl struct {
	redis *redis.Client
}

func NewLinksRepository(redisAddr string) *LinksRepositoryImpl {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	return &LinksRepositoryImpl{
		redis: rdb,
	}
}

func (r *LinksRepositoryImpl) SaveLink(ctx context.Context, url string) {
	err := r.redis.Set(ctx, url, "processed", models.TTL).Err()
	if err != nil {
		fmt.Println("redis set:", err)
		return
	}
}

func (r *LinksRepositoryImpl) CheckExists(ctx context.Context, url string) bool {
	ttl, err := r.redis.TTL(ctx, url).Result()
	if err != nil {
		fmt.Println("check TTL:", err)
		return false
	}

	if ttl <= 0 {
		fmt.Println("link dont exist:", url)
		return false
	}

	return true
}
