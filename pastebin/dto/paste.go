package dto

import "time"

type CreatePaste struct {
	Content      string `json:"content" bson:"content"`
	TTL          int64  `json:"ttl,omitempty"`
	Password     string `json:"password,omitempty"`
	AllowedEmail string `json:"allowed_email,omitempty"`
	AllowedIp    string `json:"allowed_ip,omitempty"`
	Authorized   bool   `json:"authorized,omitempty"`
}

type GetPaste struct {
	ID           string    `json:"id"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	TTL          int64     `json:"ttl,omitempty"`
	Password     string    `json:"password,omitempty"`
	AllowedEmail string    `json:"allowed_email,omitempty"`
	AllowedIp    string    `json:"allowed_ip,omitempty"`
	Authorized   bool      `json:"authorized,omitempty"`
	Deleted      bool      `json:"deleted"`
}

type PasteResponse struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type GetPasteResponse struct {
	Password string `json:"password"`
	Email    string
	IP       string
}
