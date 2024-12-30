package main

import (
	"flag"
	"log"
	"os"

	"github.com/JuDyas/GolangPractice/pastebin/internal/db"
	"github.com/JuDyas/GolangPractice/pastebin/internal/handlers"
	"github.com/JuDyas/GolangPractice/pastebin/internal/repository"
	"github.com/JuDyas/GolangPractice/pastebin/internal/routes"
	"github.com/JuDyas/GolangPractice/pastebin/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type App struct {
	DBClient     *mongo.Client
	Router       *gin.Engine
	JWTSecret    []byte
	PasteService services.PasteService
	UserService  services.UserService
}

func (app *App) Initialize(uri, port string) {
	app.DBClient = db.ConnectDatabase(uri)
	mdb := app.DBClient.Database("pastebin")
	//TODO: handle error with zap logger
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	//TODO: handle error with zap logger
	app.JWTSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(app.JWTSecret) == 0 {
		log.Fatal("JWT_SECRET env variable not set")
	}

	pasteRepo := repository.NewPasteRepository(mdb)
	userRepo := repository.NewUserRepository(mdb.Collection("users"))
	app.PasteService = services.NewPasteService(pasteRepo)
	app.UserService = services.NewUserService(userRepo)
	pasteHandler := handlers.NewPasteHandler(app.PasteService)
	app.Router = gin.Default()
	routes.SetupRoutes(app.Router, app.JWTSecret, pasteHandler, app.UserService)
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
	//TODO: handle error
	app.Router.Run(*port)
}
