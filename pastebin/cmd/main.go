package main

import (
	"log"
	"os"

	"github.com/JuDyas/GolangPractice/pastebin_new/internal/handlers"

	"github.com/JuDyas/GolangPractice/pastebin_new/internal/routes"

	"github.com/JuDyas/GolangPractice/pastebin_new/internal/repositories"
	"github.com/JuDyas/GolangPractice/pastebin_new/internal/services"

	"github.com/JuDyas/GolangPractice/pastebin_new/config"
	"github.com/JuDyas/GolangPractice/pastebin_new/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// App - application struct
type App struct {
	DBClient     *mongo.Client
	Router       *gin.Engine
	JWTSecret    []byte
	UserService  services.UserService
	PasteService services.PasteService
	AdminService services.AdminPasteService
}

func (app *App) Init(uri string) {
	//.env
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	app.JWTSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(app.JWTSecret) == 0 {
		log.Fatal("Error loading JWT_SECRET")
	}

	//Data-Base
	app.DBClient = db.ConnectDatabase(uri)
	dbName := os.Getenv("DB_NAME")
	//Auth
	userRepository := repositories.NewUserRepository(app.DBClient, dbName)
	app.UserService = services.NewUserService(userRepository)
	//Pastes
	pasteRepository := repositories.NewPasteRepository(app.DBClient, dbName)
	app.PasteService = services.NewPasteService(pasteRepository)
	pasteHandler := handlers.NewPasteHandler(app.PasteService)
	//Admin
	app.AdminService = services.NewAdminPasteService(pasteRepository)
	adminHandler := handlers.NewAdminHandler(app.AdminService)
	//Router
	app.Router = gin.Default()
	routes.SetupRoutes(app.Router, app.UserService, app.PasteService, app.JWTSecret, pasteHandler, adminHandler)
}

func main() {
	//Config
	conf, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	//Innit application
	app := App{}
	app.Init(conf.URI)
	err = app.Router.Run(conf.Port)
	if err != nil {
		//TODO: add zap logger
		log.Fatal(err)
	}
}
