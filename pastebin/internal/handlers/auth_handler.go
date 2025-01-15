package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/JuDyas/GolangPractice/pastebin_new/dto"

	"github.com/JuDyas/GolangPractice/pastebin_new/internal/auth"

	"github.com/JuDyas/GolangPractice/pastebin_new/internal/services"

	"github.com/gin-gonic/gin"
)

func Register(us services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input dto.AuthUser
		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": dto.InvalidInputErr})
			//TODO: add zap logger
			log.Println(dto.BindJsonErr, err)
			return
		}

		err := us.CreateUser(c.Request.Context(), input.Email, input.Password)
		if err != nil {
			if errors.Is(err, dto.UserExistErr) {
				c.JSON(http.StatusConflict, gin.H{"error": dto.UserExistErr.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": dto.FailedCreateUserErr})
				//TODO: add zap logger
				log.Println(dto.FailedCreateUserErr, err)
			}

			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "register success"})
	}
}

func Login(us services.UserService, jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input dto.AuthUser
		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": dto.InvalidInputErr})
			log.Println(dto.BindJsonErr, err)
		}

		user, err := us.Authenticate(c.Request.Context(), input.Email, input.Password)
		if err != nil {
			if errors.Is(err, dto.NonExistErr) {
				c.JSON(http.StatusNotFound, gin.H{"error": dto.NonExistErr.Error()})
			} else if errors.Is(err, dto.InvalidPasswordErr) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": dto.InvalidPasswordErr.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": dto.FailedLoginErr})
			}
			return
		}

		jwtToken, err := auth.GenerateJWT(user, jwtSecret)
		if err != nil {
			//TODO: add zap logger
			log.Printf("failed to generate JWT: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": dto.FailedLoginErr})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "login success",
			"token":   jwtToken})
	}
}
