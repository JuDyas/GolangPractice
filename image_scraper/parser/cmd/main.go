package main

import (
	"fmt"
	"log"
	"os"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/parser/internal/services"
)

func main() {
	fmt.Println("PARSER")
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}

	parser := services.NewWebParser(redisAddr)
	wsClient := services.NewWebSocketClient("ws://master:8080/ws/parser", parser)

	err := wsClient.Connect()
	if err != nil {
		log.Fatal("Failed to connect to masterserver:", err)
	}

	wsClient.Listen()
}
