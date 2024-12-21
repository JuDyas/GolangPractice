package main

import (
	"GolangPractice/Http_Learn/Parsing/ParsingWebSiteNew/database"
	"GolangPractice/Http_Learn/Parsing/ParsingWebSiteNew/handlers"
	"GolangPractice/Http_Learn/Parsing/ParsingWebSiteNew/tasks"
	"flag"
	"log"
	"net/http"
)

var (
	serverPort = flag.String("port", ":8080", "Port to listen on")
	cipherKey  = flag.String("key", "BanANa", "Cipher key to use")
)

func main() {
	var (
		rdb            = database.SetupRedis()
		mux            = http.NewServeMux()
		productChannel = make(chan []handlers.Product)
	)

	go tasks.InitCron(productChannel, *cipherKey, rdb)
	SetupRoutes(mux, rdb, productChannel, *cipherKey)
	err := http.ListenAndServe(*serverPort, mux)
	if err != nil {
		log.Fatalf("start server error: %s", err)
	}
}
