package main

import (
	"log"
	"os"

	"github.com/JuDyas/GolangPractice/pastebin_new/internal/routes"

	"github.com/JuDyas/GolangPractice/pastebin_new/internal/repositories"
	"github.com/JuDyas/GolangPractice/pastebin_new/internal/services"

	"github.com/JuDyas/GolangPractice/pastebin_new/config"
	"github.com/JuDyas/GolangPractice/pastebin_new/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type App struct {
	DBClient    *mongo.Client
	Router      *gin.Engine
	JWTSecret   []byte
	UserService services.UserService
}

func (app *App) Init(uri string) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	app.JWTSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(app.JWTSecret) == 0 {
		log.Fatal("Error loading JWT_SECRET")
	}

	app.DBClient = db.ConnectDatabase(uri)
	mdb := app.DBClient.Database("pastebin")
	userRepository := repositories.NewUserRepository(mdb.Collection("users"))
	app.UserService = services.NewUserService(userRepository)
	app.Router = gin.Default()
	routes.SetupRoutes(app.Router, app.UserService)
}

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	app := App{}
	app.Init(conf.URI)
	err = app.Router.Run(conf.Port)
	if err != nil {
		//TODO: add zap logger
		log.Fatal(err)
	}
}
