package main

import (
	"GolangPractice/Http_Learn/Parsing/ParsingWebSiteNew/database"
	"GolangPractice/Http_Learn/Parsing/ParsingWebSiteNew/handlers"
	"GolangPractice/Http_Learn/Parsing/ParsingWebSiteNew/tasks"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var port = flag.String("port", ":8080", "Port to listen on")

func main() {
	var (
		rdb = database.SetupRedis()
		mux = http.NewServeMux()
	)
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	key := os.Getenv("VYZHENER_KEY")
	if key == "" {
		log.Fatal("VYZHENER_KEY env variable not set")
	}

	productChannel := make(chan []handlers.Product)
	SetupRoutes(mux, rdb, productChannel)
	go tasks.InitCron(productChannel, key, rdb)
	//Graceful shutdown
	err = http.ListenAndServe(*port, mux)
	if err != nil {
		log.Fatalf("start server error: %s", err)
	}
	// Struct for handlers with object for fork with bd and cipher
}
