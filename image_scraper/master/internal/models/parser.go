package models

import (
	"time"

	"github.com/gorilla/websocket"
)

type Parser struct {
	ID           int
	Conn         *websocket.Conn
	LastPingTime time.Time
}
