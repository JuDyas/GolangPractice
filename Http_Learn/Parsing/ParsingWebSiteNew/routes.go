package main

import (
	"GolangPractice/Http_Learn/Parsing/ParsingWebSiteNew/handlers"
	"net/http"

	"github.com/go-redis/redis/v8"
)

func SetupRoutes(mux *http.ServeMux, rdb *redis.Client) {

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./html/index.html")
		return
	})

	mux.HandleFunc("/parse", handlers.ParseHtml(rdb))
	mux.HandleFunc("/getdata", handlers.GetDataFromRedis(rdb))

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))
}