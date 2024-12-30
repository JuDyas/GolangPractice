package handlers

import (
	"errors"
	"net/http"

	"github.com/JuDyas/GolangPractice/pastebin/internal/services"

	"github.com/JuDyas/GolangPractice/pastebin/internal/auth"

	"github.com/gin-gonic/gin"
)

func Register(us services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, err := us.CreateUser(c.Request.Context(), input.Email, input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.New("cannot create user")})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"user_id": userID})
	}
}

func Login(us services.UserService, jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		user, err := us.AuthenticateUser(c.Request.Context(), input.Email, input.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user data"})
			return
		}

		//TODO: вводить роль с бд
		tokenStr, err := auth.GenerateJWT(user.ID, user.Role, jwtSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenStr})
	}
}
