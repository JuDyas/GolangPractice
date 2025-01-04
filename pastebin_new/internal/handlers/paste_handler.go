package handlers

import (
	"log"
	"net/http"

	"github.com/JuDyas/GolangPractice/pastebin_new/internal/services"
	"github.com/JuDyas/GolangPractice/pastebin_new/models"
	"github.com/gin-gonic/gin"
)

type PasteHandler interface {
	CreatePaste(c *gin.Context)
}

type pasteHandlerImpl struct {
	service services.PasteService
}

func NewPasteHandler(service services.PasteService) PasteHandler {
	return pasteHandlerImpl{service: service}
}

func (h pasteHandlerImpl) CreatePaste(c *gin.Context) {
	var paste models.Paste
	if err := c.ShouldBindJSON(&paste); err != nil {
		log.Printf("bindJSON error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := h.service.CreatePaste(c.Request.Context(), &paste); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create paste error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"paste id": paste.ID})
}
