package models

import "time"

type Paste struct {
	ID           string    `json:"id" bson:"_id"`
	Text         string    `json:"text" bson:"text"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	TTL          int       `json:"time_to_live,omitempty" bson:"time_to_live"`
	Password     string    `json:"password,omitempty" bson:"password"`
	AllowedIPs   string    `json:"allowed_ips,omitempty" bson:"allowed_ips"`
	AllowedEmail string    `json:"allowed_email,omitempty" bson:"allowed_email"`
	Verified     bool      `json:"verified,omitempty" bson:"verified"`
}

type PublicPaste struct {
	ID   string `json:"id" bson:"_id"`
	Text string `json:"text" bson:"text"`
}
