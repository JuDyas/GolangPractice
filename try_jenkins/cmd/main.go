package main

import (
	"log"

	"github.com/JuDyas/GolangPractice/try_jenkins/internal/handler"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.POST("/wordcount", handler.CountWords)

	// Запускаем сервер
	log.Fatal(e.Start(":8080"))
}
