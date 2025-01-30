package main

import (
	"log"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/parseserver/routes"
	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/parseserver/services"
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

	err := app.Router.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
