package models

import "time"

type Paste struct {
	ID           string    `json:"id" bson:"_id"`
	Content      string    `json:"content" bson:"content" binding:"required"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	TTL          int64     `json:"ttl,omitempty" bson:"ttl"`
	Password     string    `json:"password,omitempty" bson:"password"`
	AllowedEmail string    `json:"allowed_email,omitempty" bson:"allowed_email"`
	AllowedIp    string    `json:"allowed_ip,omitempty" bson:"allowed_ip"`
	Authorized   bool      `json:"authorized,omitempty" bson:"authorized"`
	Deleted      bool      `bson:"deleted" json:"deleted"`
}
