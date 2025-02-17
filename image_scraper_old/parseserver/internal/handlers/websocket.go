package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/parseserver/internal/services"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebsocketHandler(ps *services.ParseService) echo.HandlerFunc {
	return func(c echo.Context) error {
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			log.Println("websocket upgrade:", err)
			return err
		}
		defer conn.Close()

		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("❌ Ошибка чтения сообщения:", err)
			return err
		}
		urlToParse := string(msg)

		log.Println("📩 Получен URL для парсинга:", urlToParse)

		imageChannel := make(chan string)
		go func() {
			err := ps.ParseImages(urlToParse, imageChannel)
			if err != nil {
				log.Println("parse images:", err)
			}
		}()

		for image := range imageChannel {
			err = conn.WriteMessage(websocket.TextMessage, []byte(image))
			if err != nil {
				log.Println("websocket write:", err)
				return err
			}
		}

		err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Parsing complete"))
		if err != nil {
			return fmt.Errorf("close websocket: %w", err)
		}

		return nil
	}
}
