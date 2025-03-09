package main

import (
	"fmt"
	"log"
	"os"

	"github.com/JuDyas/GolangPractice/image_scraper/parser/internal/handlers"
	"github.com/JuDyas/GolangPractice/image_scraper/parser/internal/models"
	"github.com/JuDyas/GolangPractice/image_scraper/parser/internal/repositories"
	"github.com/JuDyas/GolangPractice/image_scraper/parser/internal/services"
)

func main() {
	fmt.Println("PARSER")
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}

	controlChan := make(chan models.CommandType, 10)
	imageUrlChan := make(chan string)
	repo := repositories.NewLinksRepository(redisAddr)
	parser := services.NewWebParser(repo, controlChan, imageUrlChan)
	wsClient := handlers.NewWebSocketClient("ws://master:8080/ws/parser", parser, controlChan, imageUrlChan)

	err := wsClient.Connect()
	if err != nil {
		log.Fatal("Failed to connect to masterserver:", err)
	}

	wsClient.Listen()
}
