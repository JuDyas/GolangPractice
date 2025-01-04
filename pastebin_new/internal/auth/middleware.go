package auth

import "github.com/gin-gonic/gin"

func AuthoriseMiddleware(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
