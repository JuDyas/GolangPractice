package models

type CommandType string

type CommandMessage struct {
	Command CommandType `json:"command"`
	URL     string      `json:"url,omitempty"`
	Width   int         `json:"width,omitempty"`
	Height  int         `json:"height,omitempty"`
}
