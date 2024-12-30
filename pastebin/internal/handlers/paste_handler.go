package handlers

import (
	"errors"
	"net/http"

	"github.com/JuDyas/GolangPractice/pastebin/internal/services"
	"github.com/JuDyas/GolangPractice/pastebin/models"
	"github.com/gin-gonic/gin"
)

type PasteHandler struct {
	service services.PasteService
}

func NewPasteHandler(service services.PasteService) *PasteHandler {
	return &PasteHandler{service: service}
}

func (ph *PasteHandler) CreatePaste(c *gin.Context) {
	var paste models.Paste
	if err := c.ShouldBindJSON(&paste); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid data")})
		return
	}

	if err := ph.service.CreatePaste(c.Request.Context(), &paste); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.New("internal server error")})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": paste.ID})
}

func (ph *PasteHandler) GetPaste(c *gin.Context) {
	var (
		id              = c.Param("id")
		pasteExpiredErr = errors.New("paste expired")
		password        = c.Query("password")
		invalidPassword = errors.New("invalid password")
	)
	paste, err := ph.service.GetPaste(c.Request.Context(), id, password)
	if err != nil {
		if errors.Is(pasteExpiredErr, err) {
			c.JSON(http.StatusNotFound, gin.H{"error": pasteExpiredErr})
		} else if errors.Is(invalidPassword, err) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": invalidPassword})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": errors.New("paste not found")})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"paste": paste})
}
