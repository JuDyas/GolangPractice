package handlers

import (
	"GolangPractice/Http_Learn/Parsing/ParsingWebSiteNew/config"
	"GolangPractice/Http_Learn/Parsing/ParsingWebSiteNew/utils"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
)

func GetProductDetails(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are supported", http.StatusMethodNotAllowed)
			return
		}

		var requestData struct {
			DecodeKey string `json:"decodeKey"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}

		productID := strings.TrimPrefix(r.URL.Path, "/product/")
		if productID == "" {
			http.Error(w, "Product ID not specified", http.StatusBadRequest)
			return
		}

		sortOptions := []string{
			"sort-position", "sort-price", "sort-price_desc", "sort-name", "sort-name_desc", "sort-rating", "sort-rating_desc",
		}
		var product Product
		found := false

		for _, sortParam := range sortOptions {
			key := "products:" + sortParam
			jsonData, err := rdb.Get(config.Ctx, key).Result()
			if err != nil {
				continue
			}

			var products []Product
			err = json.Unmarshal([]byte(jsonData), &products)
			if err != nil {
				continue
			}

			for _, p := range products {
				if p.ID == productID {
					p.Specs = utils.DecodeData(requestData.DecodeKey, p.Specs)
					product = p
					found = true
					break
				}
			}
			if found {
				break
			}
		}

		if !found {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	}
}
