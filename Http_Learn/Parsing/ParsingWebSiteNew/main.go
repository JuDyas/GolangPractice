package main

import (
	"GolangPractice/Http_Learn/Parsing/ParsingWebSiteNew/config"
	"GolangPractice/Http_Learn/Parsing/ParsingWebSiteNew/tasks"
	"log"
	"net/http"
)

func main() {
	var (
		rdb = config.SetupRedis()
		mux = http.NewServeMux()
	)
	SetupRoutes(mux, rdb)

	fs := http.FileServer(http.Dir("./html"))
	mux.Handle("/html/", http.StripPrefix("/html", fs))
	tasks.InitCron()

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("start server error: %s", err)
	}
}
