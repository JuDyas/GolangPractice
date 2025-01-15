package models

import "time"

const (
	RoleAdmin    = 1
	RoleStandard = 2
)

type User struct {
	ID        string    `json:"id" bson:"_id"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"password" bson:"password"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	Role      int       `json:"role" bson:"role"`
}
