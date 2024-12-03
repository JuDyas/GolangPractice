package main

import (
	"GolangPractice/Http_Learn/Parsing/ParsingWebSiteNew/handlers"
	"net/http"

	"github.com/go-redis/redis/v8"
)

func SetupRoutes(mux *http.ServeMux, rdb *redis.Client, productChannel chan []handlers.Product, key string) {

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./html/index.html")
		return
	})

	mux.HandleFunc("/parse", handlers.ParseHtml(productChannel))
	mux.HandleFunc("/products", handlers.GetProducts(rdb))
	mux.HandleFunc("/products/", handlers.GetProduct(rdb, key))
}
