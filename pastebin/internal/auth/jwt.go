package auth

import "github.com/golang-jwt/jwt/v5"

func GenerateJWT(userID, role string, jwtSecret []byte) (string, error) {
	var claims = jwt.MapClaims{
		"userID": userID,
		"role":   role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
