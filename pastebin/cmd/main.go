package main

import (
	"flag"
	"log"
	"os"

	"github.com/JuDyas/GolangPractice/pastebin/internal/routes"

	"github.com/joho/godotenv"

	"github.com/JuDyas/GolangPractice/pastebin/internal/db"
	"github.com/gin-gonic/gin"
)

var (
	jwtSecret []byte
	uri       = flag.String("uri", "mongodb://localhost:27017", "mongo database URI")
	port      = flag.String("port", ":8080", "port to listen on")
)

func main() {
	db.ConnectDatabase(*uri)
	//TODO: handle error with zap logger
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	//TODO: handle error with zap logger
	if len(jwtSecret) == 0 {
		log.Fatal("JWT_SECRET env variable not set")
	}

	r := gin.Default()
	routes.SetupRoutes(r, jwtSecret, *port)
}
