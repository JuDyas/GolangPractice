package tasks

import (
	"GolangPractice/Http_Learn/Parsing/ParsingWebSiteNew/handlers"
	"GolangPractice/Http_Learn/Parsing/ParsingWebSiteNew/utils"
	"GolangPractice/Http_Learn/Parsing/vyzhenercipher"
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
			"ID", id,
			"Name", p.Name,
			"CPU", vyzhenercipher.Encode(p.Specs.Cpu, key),
			"GPU", vyzhenercipher.Encode(p.Specs.Gpu, key),
			"DisplaySize", vyzhenercipher.Encode(p.Specs.DisplaySize, key),
			"DisplayResolution", vyzhenercipher.Encode(p.Specs.DisplayResolution, key),
			"RAM", vyzhenercipher.Encode(p.Specs.Ram, key),
			"HardDrives", vyzhenercipher.Encode(p.Specs.HardDrives, key),
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
