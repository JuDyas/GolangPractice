package auth

import (
	"log"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT - generate jwt token and signe it with secret
func GenerateJWT(userID, role string, jwtSecret []byte) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//TODO: add zap logger
	log.Printf("created token for user: %v", userID)
	return token.SignedString(jwtSecret)
}
