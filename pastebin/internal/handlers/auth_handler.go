package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/JuDyas/GolangPractice/pastebin_new/internal/auth"

	"github.com/JuDyas/GolangPractice/pastebin_new/internal/services"

	"github.com/gin-gonic/gin"
)

func Register(us services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=8,max=20"`
		}

		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
			//TODO: add zap logger
			log.Printf("bindJson error: %v", err)
			return
		}

		err := us.CreateUser(c.Request.Context(), input.Email, input.Password)
		if err != nil {
			if errors.Is(err, fmt.Errorf("user already exists")) {
				c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
				//TODO: add zap logger
				log.Printf("failed to create user: %v", err)
			}

			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "register success"})
	}
}

func Login(us services.UserService, jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=8,max=20"`
		}

		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			log.Printf("bindJson error: %v", err)
		}

		user, err := us.Authenticate(c.Request.Context(), input.Email, input.Password)
		if err != nil {
			if errors.Is(err, fmt.Errorf("user not found")) {
				c.JSON(http.StatusNotFound, gin.H{"error": "user does not exists"})
			} else if errors.Is(err, fmt.Errorf("password incorrect")) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
			}
			return
		}

		jwtToken, err := auth.GenerateJWT(user, jwtSecret)
		if err != nil {
			//TODO: add zap logger
			log.Printf("failed to generate JWT: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "login success",
			"token":   jwtToken})
	}
}
