package models

import "time"

type Paste struct {
	ID         string    `json:"id" bson:"_id"`
	Text       string    `json:"text" bson:"text"`
	ExpiresAt  time.Time `json:"expires_at" bson:"expires_at"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	Password   string    `json:"password,omitempty" bson:"password"`
	AllowedIPs []string  `json:"allowed_ips,omitempty" bson:"allowed_ips"`
	OwnerEmail string    `json:"owner_email,omitempty" bson:"owner_email"`
}
