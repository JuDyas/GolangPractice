package handlers

import (
	"net/http"

	"github.com/JuDyas/GolangPractice/pastebin/internal/auth"

	"github.com/JuDyas/GolangPractice/pastebin/internal/servises"
	"github.com/JuDyas/GolangPractice/pastebin/models"
	"github.com/gin-gonic/gin"
)

func Register(us *servises.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.User
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, err := us.CreateUser(input.Email, input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"user": userID})

	}
}

func Login(us *servises.UserService, jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.User
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := us.CheckPassword(input.Email, input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		//TODO: вводить роль с бд
		tokenStr, err := auth.GenerateJWT(user.ID, "standard", jwtSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.SetCookie("token", tokenStr, 259000, "/", "localhost", false, true)
		c.JSON(http.StatusOK, gin.H{"massage": "logged in"})
	}
}
