package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/JuDyas/GolangPractice/pastebin_new/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthoriseMiddleware(jwtSecret []byte, requiredRole ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			//TODO: add zap logger
			log.Printf("validate token error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			log.Printf("get token claims error")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		}

		email, ok := claims["mail"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token payload"})
			c.Abort()
			return
		}

		c.Set("email", email)
		if len(requiredRole) > 0 {
			role, ok := claims["role"].(string)
			if !ok || !contains(requiredRole, role) {
				c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
				c.Abort()
				return
			}
			c.Set("role", role)
		}

		c.Next()
	}
}

func contains(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func PasteMiddleware(ps services.PasteService, jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		pasteID := c.Param("id")
		paste, err := ps.GetPasteByID(c.Request.Context(), pasteID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
			c.Abort()
			return
		}

		c.Set("paste", paste)
		if paste.AllowedEmail == "" && paste.Authorized == false {
			c.Next()
			return
		}

		authorise := AuthoriseMiddleware(jwtSecret)
		authorise(c)
		if c.IsAborted() {
			return
		}

		c.Next()
	}
}
