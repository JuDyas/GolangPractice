package handlers

import (
	"fmt"
	"net/http"
	"sync"

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
	imgWebChan   chan string
	imgLocalChan chan string
	cmdChan      chan string
	service      *services.ImageServiceImpl
}

func NewWebSocket(s *services.ImageServiceImpl) *WebSocketImpl {
	ws := &WebSocketImpl{
		imgWebChan:   make(chan string),
		imgLocalChan: make(chan string),
		cmdChan:      make(chan string),
		service:      s,
	}

	go ws.processImages()
	return ws
}

func (ws *WebSocketImpl) ParserHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			fmt.Println("ws parser conn:", err)
			return err
		}
		defer conn.Close()

		var wg = new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					fmt.Println("read message from parser:", err)
					return
				}

				fmt.Println("get from parser:", string(msg))
				ws.imgWebChan <- string(msg)
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for cmd := range ws.cmdChan {
				err := conn.WriteMessage(websocket.TextMessage, []byte(cmd))
				if err != nil {
					fmt.Println("write message to parser:", err)
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

func (ws *WebSocketImpl) processImages() {
	for img := range ws.imgWebChan {
		url, err := ws.service.ProcessImage(img)
		if err != nil {
			fmt.Println("process image err:", err)
			continue
		}
		fmt.Println("image was processed:", url)
		ws.imgLocalChan <- url
	}
}

func (ws *WebSocketImpl) ClientHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			fmt.Println("ws client conn:", err)
			return err
		}
		defer conn.Close()

		wg := new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					fmt.Println("read message from client:", err)
					return
				}

				fmt.Println("get cmd from client:", string(msg))
				ws.cmdChan <- string(msg)
			}
		}()

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
