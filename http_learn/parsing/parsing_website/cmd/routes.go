package main

import (
	handlers2 "GolangPractice/http_learn/parsing/parsing_website/internal/handlers"
	"net/http"

	"github.com/go-redis/redis/v8"
)

func SetupRoutes(mux *http.ServeMux, rdb *redis.Client, productChannel chan []handlers2.Product, key string) {

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./html/index.html")
		return
	})

	mux.HandleFunc("/parse", handlers2.ParseHtml(productChannel))
	mux.HandleFunc("/products", handlers2.GetProducts(rdb))
	mux.HandleFunc("/products/", handlers2.GetProduct(rdb, key))
}
