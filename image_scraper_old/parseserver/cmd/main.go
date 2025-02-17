package main

import (
	"fmt"
	"log"
	"os"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/parseserver/internal/services"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/parseserver/routes"
	"github.com/labstack/echo/v4"
)

type App struct {
	Router       *echo.Echo
	ParseService *services.ParseService
}

func (a *App) Initialize() {
	a.Router = echo.New()
	a.ParseService = &services.ParseService{
		Parser: &services.CollyParser{},
	}

	routes.SetupRoutes(a.Router, a.ParseService)
}

func main() {
	app := App{}
	app.Initialize()

	serverPort := os.Getenv("PARSER_SERVER_PORT")
	if serverPort == "" {
		serverPort = "8081"
		fmt.Println("parser port not found in environment variable PARSER_SERVER_PORT, use:", serverPort)
	} else {
		fmt.Println("parser port found:", serverPort)
	}

	err := app.Router.Start(":" + serverPort)
	if err != nil {
		log.Fatal(err)
	}
}
