package main

import (
	"fmt"
	"log"
	"os"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/masterserver/internal/repositories"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/masterserver/routes"
	"github.com/labstack/echo/v4"
)

type App struct {
	Router *echo.Echo
}

func (a *App) Init() {
	a.Router = echo.New()
	repositories.InitDB()
	routes.SetupRoutes(a.Router)
}

func main() {
	app := App{}
	app.Init()
	port := os.Getenv("MASTER_SERVER_PORT")
	fmt.Println("ПОРТ: ", port)
	log.Fatal(app.Router.Start(":" + port))
}
