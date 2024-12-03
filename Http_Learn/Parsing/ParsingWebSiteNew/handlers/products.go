package handlers

import (
	"GolangPractice/Http_Learn/Parsing/ParsingWebSiteNew/utils"
	"GolangPractice/Http_Learn/Parsing/vyzhenercipher"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-redis/redis/v8"
)

func GetProducts(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx      = r.Context()
			prodName []string
		)
		keys, err := rdb.Keys(ctx, "product:*").Result()
		if err != nil {
			http.Error(w, "Get data from db error 1", http.StatusInternalServerError)
			return
		}

		for _, key := range keys {
			data, err := rdb.HGetAll(ctx, key).Result()
			if err != nil {
				http.Error(w, "Get data from db error 2", http.StatusInternalServerError)
				return
			}

			if name, exists := data["name"]; exists {
				prodName = append(prodName, name)
				prodName = append(prodName, utils.HashMD5(name))
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prodName)
	}
}

func GetProduct(rdb *redis.Client, vKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx = r.Context()
			id  = r.URL.Path[len("/products/"):]
			key = r.URL.Query().Get("key")
		)
		if id == "" {
			http.Error(w, "Product ID is required", http.StatusBadRequest)
			return
		}

		err := checkKey(key, vKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data, err := rdb.HGetAll(ctx, "product:"+id).Result()
		if err == redis.Nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Failed to fetch product details", http.StatusInternalServerError)
		}

		decodedData := map[string]string{}
		for k, v := range data {
			if k != "name" {
				decodedData[k] = vyzhenercipher.Decode(v, key)
			} else {
				decodedData[k] = v
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(decodedData)
	}
}

func checkKey(key, vKey string) error {
	if key == "" {
		return errors.New("key is required")
	} else if key != vKey {
		return errors.New("key is wrong")
	} else {
		return nil
	}
}
