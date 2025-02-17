package main

import (
	"fmt"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/master/internal/db"
	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/master/internal/repositories"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/master/internal/handlers"
	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/master/internal/services"
	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/master/routes"

	"github.com/labstack/echo/v4"
)

//type CommandType string
//
//const (
//	CmdStart    CommandType = "start"
//	CmdStop     CommandType = "stop"
//	CmsPause    CommandType = "pause"
//	CmdContinue CommandType = "continue"
//)
//
//type CommandMessage struct {
//	Command CommandType `json:"command"`
//	URL     string      `json:"url,omitempty"`
//}
//
//var upgrader = websocket.Upgrader{
//	CheckOrigin: func(r *http.Request) bool {
//		return true
//	},
//}
//
//var clients = make(map[*websocket.Conn]bool) // Храним всех подключенных клиентов
//
//func websocketParser(c echo.Context) error {
//	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
//	if err != nil {
//		log.Println("websocket conn:", err)
//		return err
//	}
//	defer conn.Close()
//
//	clients[conn] = true
//	sendCommand(CmdStart, "https://www.moyo.ua/")
//	log.Println("clients:", len(clients))
//
//	for {
//		_, _, err := conn.ReadMessage() // Просто держим соединение
//		if err != nil {
//			log.Println("Error reading message:", err)
//			delete(clients, conn)
//			break
//		}
//	}
//	return nil
//}
//
//func sendCommand(command CommandType, url string) {
//	fmt.Println("SEND COMM")
//	msg := CommandMessage{
//		Command: command,
//		URL:     url,
//	}
//	data, err := json.Marshal(msg)
//	if err != nil {
//		log.Println("JSON Marshal error:", err)
//		return
//	}
//
//	for client := range clients {
//		err := client.WriteMessage(websocket.TextMessage, data)
//		if err != nil {
//			log.Println("Error sending command:", err)
//			client.Close()
//			delete(clients, client)
//		}
//	}
//}

func main() {
	fmt.Println("MASTER")
	e := echo.New()
	postgres := db.InitPostgres()
	repo := repositories.NewImageRepository(postgres)
	serv := services.NewImageService(repo)
	webs := handlers.NewWebSocket(serv)
	routes.SetupRoutes(e, webs)
	e.Start(":8080")
}
