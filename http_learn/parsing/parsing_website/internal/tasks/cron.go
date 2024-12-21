package tasks

import (
	"GolangPractice/http_learn/parsing/parsing_website/internal/handlers"
	"GolangPractice/http_learn/parsing/parsing_website/pkg/vyzhenercipher"
	"GolangPractice/http_learn/parsing/parsing_website/utils"
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

const zsetKey = "products:all_keys"

func InitCron(productChannel chan []handlers.Product, key string, rdb *redis.Client) {
	for prod := range productChannel {
		processProducts(prod, rdb, key)
	}
}

func processProducts(prod []handlers.Product, rdb *redis.Client, key string) {
	var (
		ctx = context.Background()
	)

	for _, p := range prod {
		id := utils.HashMD5(p.Name)
		keyDB := "product:" + id

		err := rdb.HSet(ctx, keyDB, []string{
			"Id", id,
			"Name", p.Name,
			"Cpu", vyzhenercipher.Encode(p.Details.Cpu, key),
			"Gpu", vyzhenercipher.Encode(p.Details.Gpu, key),
			"Display_size", vyzhenercipher.Encode(p.Details.DisplaySize, key),
			"Display_resolution", vyzhenercipher.Encode(p.Details.DisplayResolution, key),
			"Ram", vyzhenercipher.Encode(p.Details.Ram, key),
			"Hard_drives", vyzhenercipher.Encode(p.Details.HardDrives, key),
		}).Err()
		if err != nil {
			continue
		}

		err = rdb.ZAdd(ctx, zsetKey, &redis.Z{
			Score:  float64(time.Now().Unix()),
			Member: keyDB,
		}).Err()

		if err != nil {
			continue
		}
	}

}
