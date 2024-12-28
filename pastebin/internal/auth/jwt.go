package auth

import "github.com/golang-jwt/jwt/v5"

func GenerateJWT(userID, role string, jwtSecret []byte) (string, error) {
	//TODO: Add exp (ttl for jwt token)
	var claims = jwt.MapClaims{
		"userID": userID,
		"role":   role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
