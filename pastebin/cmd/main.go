package main

import (
	"flag"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/JuDyas/GolangPractice/pastebin/internal/routes"

	"github.com/joho/godotenv"

	"github.com/JuDyas/GolangPractice/pastebin/internal/db"
	"github.com/gin-gonic/gin"
)

type App struct {
	DBClient  *mongo.Client
	Router    *gin.Engine
	JWTSecret []byte
}

func (app *App) Initialize(uri, port string) {
	app.DBClient = db.ConnectDatabase(uri)
	//TODO: handle error with zap logger
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	//TODO: handle error with zap logger
	app.JWTSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(app.JWTSecret) == 0 {
		log.Fatal("JWT_SECRET env variable not set")
	}

	app.Router = gin.Default()
	routes.SetupRoutes(app.Router, app.JWTSecret, port)
}

func main() {
	// TODO: Перенести флаги в конфиг
	var (
		uri  = flag.String("uri", "mongodb://localhost:27017", "mongo database URI")
		port = flag.String("port", ":8080", "port to listen on")
		app  = &App{}
	)
	flag.Parse()
	app.Initialize(*uri, *port)
	app.Router.Run(*port)
}
