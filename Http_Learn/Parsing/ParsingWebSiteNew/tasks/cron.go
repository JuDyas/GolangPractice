package tasks

import (
	"GolangPractice/Http_Learn/Parsing/ParsingWebSiteNew/handlers"
	"GolangPractice/Http_Learn/Parsing/vyzhenercipher"
	"log"

	"github.com/go-redis/redis/v8"
)

func InitCron(productChannel chan []handlers.Product, key string, rdb *redis.Client) {
	for products := range productChannel {
		processProducts(products, rdb, key)
	}
}

func processProducts(prod []handlers.Product, rdb *redis.Client, key string) {
	for i, product := range prod {
		encryptSpecs := vyzhenercipher.Encode(product.Specs, key)
		err := saveInDB(rdb, product.Name, encryptSpecs, i)
		if err != nil {
			log.Fatal("Save data in db error", err)
		}
	}
}

func saveInDB(rdb *redis.Client, productName, encryptSpecs string, id int) error {
	//save data func
	return nil
}
