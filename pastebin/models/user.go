package models

import "time"

type User struct {
	ID        string    `json:"id" bson:"_id"`
	Email     string    `json:"email" binding:"required,email" bson:"email"`
	Password  string    `json:"password" binding:"required,min=6" bson:"password"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	Role      string    `json:"role" bson:"role"`
}
