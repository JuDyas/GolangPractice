package routes

import (
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	v1 := e.Group("api/v1")
	v1.GET("/", func(c echo.Context) error {
		return nil
	})
}
