package dto

import (
	"errors"
	"time"
)

var (
	UserExistErr        = errors.New("user already exists")
	FailedCreateUserErr = errors.New("failed to create user")
	InvalidInputErr     = errors.New("invalid input")
	BindJsonErr         = errors.New("bindJson")
	NonExistErr         = errors.New("user do not exist")
	NotFoundErr         = errors.New("user not found")
	InvalidPasswordErr  = errors.New("password is incorrect")
	FailedLoginErr      = errors.New("failed to login")
)

type User struct {
	ID        string
	Email     string
	Password  string
	CreatedAt time.Time
	Role      int
}
type AuthUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
}
