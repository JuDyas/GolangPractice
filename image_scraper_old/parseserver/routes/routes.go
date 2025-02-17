package routes

import (
	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/parseserver/internal/handlers"
	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/parseserver/internal/services"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, ps *services.ParseService) {
	e.GET("/ws", handlers.WebsocketHandler(ps))
}
