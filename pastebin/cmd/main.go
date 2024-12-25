package main

import (
	"flag"

	"github.com/JuDyas/GolangPractice/pastebin/internal/db"
	"github.com/JuDyas/GolangPractice/pastebin/internal/handlers"
	"github.com/gin-gonic/gin"
)

var (
	uri  = flag.String("uri", "mongodb://localhost:27017", "mongo database URI")
	port = flag.String("port", ":8080", "port to listen on")
)

func main() {
	db.ConnectDatabase(*uri)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "Pastebin API is running.",
			"info":   "go to /pastes to create past",
		})
	})

	//TODO: Add other routes
	r.POST("/pastes", handlers.CreatePaste)
	r.GET("/pastes/:id", handlers.GetPaste)
	//TODO: handle error with zap logger
	r.Run(*port)
}
