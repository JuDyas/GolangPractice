package handlers

import (
	"GolangPractice/Http_Learn/Parsing/ParsingWebSiteNew/config"
	"GolangPractice/Http_Learn/Parsing/ParsingWebSiteNew/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
)

func GetDataFromRedis(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sortParam := r.URL.Query().Get("sort")
		if sortParam == "" {
			sortParam = "sort-position"
		}

		key := "products:" + sortParam
		jsonData, err := rdb.Get(config.Ctx, key).Result()
		if err != nil {
			log.Printf("error getting data from Redis %s: %v", key, err)
			http.Error(w, "Error retrieving data", http.StatusInternalServerError)
			return
		}

		var products []Product
		err = json.Unmarshal([]byte(jsonData), &products)
		if err != nil {
			log.Printf("error unmarshalling data %s: %v", key, err)
			http.Error(w, "Error processing data", http.StatusInternalServerError)
			return
		}

		for i, product := range products {
			products[i].Specs = utils.DecodeData(utils.Key, product.Specs)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(products)
		if err != nil {
			log.Printf("error encoding data: %v", err)
			return
		}
	}
}
