package auth

import "github.com/gin-gonic/gin"

func AuthorizeMiddleware(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
