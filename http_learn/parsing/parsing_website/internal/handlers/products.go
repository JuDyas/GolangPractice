package handlers

import (
	"GolangPractice/http_learn/parsing/parsing_website/pkg/vyzhenercipher"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type productResponse struct {
	Products []ProductDTO `json:"products"`
}

type ProductDTO struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Cpu                string `json:"cpu,omitempty" cipher:"true"`
	RAM                string `json:"ram,omitempty" cipher:"true"`
	Display_size       string `json:"display_size,omitempty" cipher:"true"`
	Display_resolution string `json:"display_resolution,omitempty" cipher:"true"`
	Hard_drives        string `json:"hard_drive,omitempty" cipher:"true"`
	GPU                string `json:"gpu,omitempty" cipher:"true"`
}

func GetProducts(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx               = r.Context()
			productNames      []ProductDTO
			startPag, stopPag int64
			page              = 1
			pageSize          = 10
		)

		if p := r.URL.Query().Get("page"); p != "" {
			page, _ = strconv.Atoi(p)
		}

		if ps := r.URL.Query().Get("pageSize"); ps != "" {
			pageSize, _ = strconv.Atoi(ps)
		}

		startPag = int64((page - 1) * pageSize)
		stopPag = startPag + int64(pageSize) - 1

		keys, err := rdb.ZRange(ctx, "products:all_keys", startPag, stopPag).Result()
		if err != nil {
			http.Error(w, "cant fetch keys", http.StatusInternalServerError)
			return
		}

		for _, key := range keys {
			data, err := rdb.HGetAll(ctx, key).Result()
			if err != nil {
				http.Error(w, "can not fetch data", http.StatusInternalServerError)
				return
			}

			if name, exists := data["Name"]; exists {
				productNames = append(productNames, ProductDTO{Name: name, ID: data["ID"]})
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&productResponse{Products: productNames})
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
		if errors.Is(err, redis.Nil) {
			http.Error(w, "product not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "failed to fetch product details", http.StatusInternalServerError)
		}

		product := parseProduct(data, vKey)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	}
}

func parseProduct(data map[string]string, cipherKey string) ProductDTO {
	var product ProductDTO
	productVal := reflect.ValueOf(&product).Elem()
	productType := productVal.Type()
	for i := 0; i < productType.NumField(); i++ {
		field := productType.Field(i)
		tag := field.Tag.Get("cipher")
		fieldValue := productVal.FieldByName(field.Name)

		if !fieldValue.CanSet() {
			continue
		}

		if val, ok := data[field.Name]; ok {
			if tag == "true" {
				fieldValue.SetString(vyzhenercipher.Decode(val, cipherKey))
			} else {
				fieldValue.SetString(val)
			}
		}
	}
	return product
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
