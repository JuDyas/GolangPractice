package tasks

import (
	"GolangPractice/Http_Learn/Parsing/ParsingWebSiteNew/handlers"
	"GolangPractice/Http_Learn/Parsing/ParsingWebSiteNew/utils"
	"GolangPractice/Http_Learn/Parsing/vyzhenercipher"
	"context"

	"github.com/go-redis/redis/v8"
)

func InitCron(productChannel chan []handlers.Product, key string, rdb *redis.Client) {
	for prod := range productChannel {
		processProducts(prod, rdb, key)
	}
}

func processProducts(prod []handlers.Product, rdb *redis.Client, key string) {
	ctx := context.Background()
	for _, p := range prod {
		keyDB := "product:" + utils.HashMD5(p.Name)
		_ = rdb.HSet(ctx, keyDB, []string{
			"name", p.Name,
			"cpu", vyzhenercipher.Encode(p.Specs.Cpu, key),
			"gpu", vyzhenercipher.Encode(p.Specs.Gpu, key),
			"displaySize", vyzhenercipher.Encode(p.Specs.DisplaySize, key),
			"displayResolution", vyzhenercipher.Encode(p.Specs.DisplayResolution, key),
			"ram", vyzhenercipher.Encode(p.Specs.Ram, key),
			"hardDrive", vyzhenercipher.Encode(p.Specs.HardDrives, key),
		}).Err()
	}
}
