package main

import (
	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/parseserver/routes"
	"github.com/labstack/echo/v4"
)

type App struct {
	Router *echo.Echo
}

func (a *App) Init() {
	a.Router = echo.New()
	routes.SetupRoutes(a.Router)
}

func main() {
	app := App{}
	app.Init()
	app.Router.Start("")
}
