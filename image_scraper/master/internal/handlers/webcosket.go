package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/master/internal/models"

	"github.com/JuDyas/GolangPractice/pastebin_new/image_scraper/master/internal/services"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocket interface {
	ParserHandler() echo.HandlerFunc
	ClientHandler() echo.HandlerFunc
}

type WebSocketImpl struct {
	cmdMsg       models.CommandMessage
	imgWebChan   chan string
	imgLocalChan chan string
	cmdChan      chan string
	service      *services.ImageServiceImpl
	parsers      map[int]*models.Parser
	nextParserID int
	pingInterval time.Duration
	ctx          context.Context
}

// NewWebSocket - Create new WebSocket interface
func NewWebSocket(s *services.ImageServiceImpl) *WebSocketImpl {
	ctx, cancel := context.WithCancel(context.Background())

	ws := &WebSocketImpl{
		imgWebChan:   make(chan string),
		imgLocalChan: make(chan string),
		cmdChan:      make(chan string),
		service:      s,
		parsers:      make(map[int]*models.Parser),
		nextParserID: 1,
		pingInterval: 30 * time.Second,
		ctx:          ctx,
	}

	go ws.processImages(ctx)
	go ws.monitorParsers(ctx)
	go ws.pingParsers(ctx, cancel)
	return ws
}

func (ws *WebSocketImpl) addParser(conn *websocket.Conn) int {
	id := ws.nextParserID
	ws.parsers[id] = &models.Parser{
		ID:           id,
		Conn:         conn,
		LastPingTime: time.Now(),
	}

	ws.nextParserID++
	fmt.Printf("parser %d connected\n", id)
	return id
}

func (ws *WebSocketImpl) removeParser(id int) {
	if _, ok := ws.parsers[id]; ok {
		delete(ws.parsers, id)
		fmt.Printf("parser %d was removed\n", id)
	}
}

func (ws *WebSocketImpl) pingParsers(ctx context.Context, cancel context.CancelFunc) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for id, parser := range ws.parsers {
				err := parser.Conn.WriteMessage(websocket.PingMessage, []byte("ping"))
				if err != nil {
					fmt.Printf("failed to ping parser %d: %v\n", id, err)
					ws.removeParser(id)
				}
			}
		}
	}
}

// ParserHandler - Connect to parser and write/read messages
func (ws *WebSocketImpl) ParserHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			fmt.Println("ws parser conn:", err)
			return err
		}
		defer conn.Close()

		id := ws.addParser(conn)
		defer ws.removeParser(id)

		var wg = new(sync.WaitGroup)
		// read messages from parser
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ws.ctx.Done():
					return
				default:
					_, msg, err := conn.ReadMessage()
					if err != nil {
						fmt.Println("read message from parser:", err)
						return
					}

					fmt.Println("get from parser:", string(msg))
					ws.imgWebChan <- string(msg)
				}
			}
		}()

		// read commands from client and send is to parser
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ws.ctx.Done():
					return
				case cmd, ok := <-ws.cmdChan:
					if !ok {
						return
					}

					err := conn.WriteMessage(websocket.TextMessage, []byte(cmd))
					if err != nil {
						fmt.Println("write message to parser:", err)
						return
					}
				}
			}
		}()

		wg.Wait()
		err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Parsing was stopped"))
		if err != nil {
			return fmt.Errorf("close websocket: %w", err)
		}

		return nil
	}
}

// processImages - Process images from url in channel
func (ws *WebSocketImpl) processImages(ctx context.Context) {
	for {
		select {
		case img := <-ws.imgWebChan:
			url, err := ws.service.ProcessImage(img, ws.cmdMsg.Width, ws.cmdMsg.Height)
			if err != nil {
				fmt.Println("process image err:", err)
				continue
			}

			fmt.Println("Image was processed:", url)
			ws.imgLocalChan <- url
		case <-ctx.Done():
			fmt.Println("processImages goroutine stopped")
			return
		}
	}
}

func (ws *WebSocketImpl) monitorParsers(ctx context.Context) {
	for {
		select {
		case <-ctx.Done(): // Проверяем, завершён ли контекст
			fmt.Println("monitorParsers goroutine stopped")
			return
		default:
			// Логика мониторинга парсеров
			time.Sleep(5 * time.Second) // Пример, можно проверить активность парсеров
		}
	}
}

// ClientHandler - Connect to client and write/read messages
func (ws *WebSocketImpl) ClientHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			fmt.Println("ws client conn:", err)
			return err
		}
		defer conn.Close()

		wg := new(sync.WaitGroup)
		// read command messages and send is to parser channel
		wg.Add(1)
		go func() {
			//TODO close channels
			defer wg.Done()
			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					fmt.Println("read message from client:", err)
					return
				}

				err = json.Unmarshal(msg, &ws.cmdMsg)
				if err != nil {
					log.Println("error unmarshalling message:", err)
					continue
				}

				fmt.Println("get cmd from client:", string(msg))
				ws.cmdChan <- string(msg)
			}
		}()

		// write image url from parser channel
		wg.Add(1)
		go func() {
			defer wg.Done()
			for img := range ws.imgLocalChan {
				err := conn.WriteMessage(websocket.TextMessage, []byte(img))
				if err != nil {
					fmt.Println("write message to client:", err)
					return
				}
			}
		}()

		wg.Wait()
		err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Parsing complete"))
		if err != nil {
			return fmt.Errorf("close websocket: %w", err)
		}

		return nil
	}
}
