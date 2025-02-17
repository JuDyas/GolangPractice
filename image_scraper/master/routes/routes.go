package routes

import (
	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/master/internal/handlers"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, ws handlers.WebSocket) {
	e.GET("/ws/parser", ws.ParserHandler())
	e.GET("/ws/client", ws.ClientHandler())
}
