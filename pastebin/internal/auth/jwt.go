package auth

import (
	"log"

	"github.com/JuDyas/GolangPractice/pastebin_new/models"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT - generate jwt token and signe it with secret
func GenerateJWT(user *models.User, jwtSecret []byte) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"mail": user.Email,
		"role": user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//TODO: add zap logger
	log.Printf("created token for user: %v", user.ID)
	return token.SignedString(jwtSecret)
}
