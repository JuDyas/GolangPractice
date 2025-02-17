package services

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type CommandType string

const (
	CmdStart    CommandType = "start"
	CmdStop     CommandType = "stop"
	CmsPause    CommandType = "pause"
	CmdContinue CommandType = "continue"
)

type CommandMessage struct {
	Command CommandType `json:"command"`
	URL     string      `json:"url,omitempty"`
}

type WebSocketClientImpl struct {
	url    string
	conn   *websocket.Conn
	parser Parser
}

func NewWebSocketClient(url string, parser Parser) *WebSocketClientImpl {
	return &WebSocketClientImpl{
		url:    url,
		parser: parser,
	}
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
			log.Println("Error reading message:", err)
			break
		}

		var cmdMsg CommandMessage
		err = json.Unmarshal(message, &cmdMsg)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		log.Println("Received command:", cmdMsg.Command)

		switch cmdMsg.Command {
		case CmdStart:
			if cmdMsg.URL == "" {
				log.Println("Start command received without URL")
				continue
			}
			log.Println("Starting parser for URL:", cmdMsg.URL)
			w.parser.Start(cmdMsg.URL)
		case CmdStop:
			w.parser.Stop()
		case CmsPause:
			w.parser.Pause()
		case CmdContinue:
			w.parser.Continue()
		default:
			log.Println("Unknown command:", cmdMsg.Command)
		}
	}
}

func (w *WebSocketClientImpl) Close() {
	if w.conn != nil {
		err := w.conn.Close()
		if err != nil {
			log.Println("Error closing websocket connection")
		}
	}
}
