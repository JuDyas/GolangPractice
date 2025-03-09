package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/JuDyas/GolangPractice/image_scraper/parser/internal/models"
	"github.com/JuDyas/GolangPractice/image_scraper/parser/internal/services"

	"github.com/gorilla/websocket"
)

type WebSocketClientImpl struct {
	url          string
	conn         *websocket.Conn
	parser       services.Parser
	controlChan  chan<- models.CommandType
	imageUrlChan <-chan string
}

func NewWebSocketClient(url string, parser services.Parser, controlChan chan<- models.CommandType, imageUrlChan <-chan string) *WebSocketClientImpl {
	ws := &WebSocketClientImpl{
		url:          url,
		parser:       parser,
		controlChan:  controlChan,
		imageUrlChan: imageUrlChan,
	}

	go ws.writeMessage()
	return ws
}

func (w *WebSocketClientImpl) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(w.url, nil)
	if err != nil {
		return fmt.Errorf("websocket dial error: %s", err)
	}

	w.conn = conn
	log.Println("Connected to master server")
	return nil
}

func (w *WebSocketClientImpl) Listen() {
	defer func() {
		if w.conn != nil {
			err := w.conn.Close()
			if err != nil {
				log.Println("Error closing websocket connection")
			}
		}
	}()

	for {
		_, message, err := w.conn.ReadMessage()
		if err != nil {
			log.Println("error reading message:", err)
			break
		}

		var cmdMsg models.CommandMessage
		err = json.Unmarshal(message, &cmdMsg)
		if err != nil {
			log.Println("error unmarshalling message:", err)
			continue
		}

		log.Println("Received command:", cmdMsg.Command)

		switch cmdMsg.Command {
		case models.CmdStart:
			if cmdMsg.URL == "" {
				log.Println("Start command received without URL")
				continue
			}
			log.Println("Starting parser for URL:", cmdMsg.URL)
			w.parser.Start(cmdMsg.URL)
		default:
			w.controlChan <- cmdMsg.Command
		}
	}
}

func (w *WebSocketClientImpl) writeMessage() {
	for imgSrc := range w.imageUrlChan {
		err := w.conn.WriteMessage(websocket.TextMessage, []byte(imgSrc))
		if err != nil {
			fmt.Println("write message to parser:", err)
			return
		}
	}
}

func (w *WebSocketClientImpl) Close() {
	if w.conn != nil {
		err := w.conn.Close()
		if err != nil {
			log.Println("error closing websocket connection")
		}
	}
}
