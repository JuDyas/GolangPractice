package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/master/internal/db"
	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/master/internal/repositories"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/master/internal/handlers"
	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/master/internal/services"
	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/master/routes"

	"github.com/labstack/echo/v4"
)

type App struct {
	router *echo.Echo
	db     *sql.DB
}

func (app *App) Initialize() error {
	var err error
	app.router = echo.New()
	app.db, err = db.InitPostgres()
	if err != nil {
		return err
	}

	var (
		repo = repositories.NewImageRepository(app.db)
		serv = services.NewImageService(repo)
		webs = handlers.NewWebSocket(serv)
	)

	routes.SetupRoutes(app.router, webs)

	return nil
}

func main() {
	app := App{}
	if err := app.Initialize(); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("MASTER_SERVER_PORT")
	app.router.Start(":" + port)
}
