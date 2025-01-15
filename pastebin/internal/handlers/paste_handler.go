package handlers

import (
	"log"
	"net/http"

	"github.com/JuDyas/GolangPractice/pastebin_new/dto"

	"github.com/JuDyas/GolangPractice/pastebin_new/internal/services"
	"github.com/gin-gonic/gin"
)

type PasteHandler interface {
	CreatePaste(c *gin.Context)
	GetPaste(c *gin.Context)
}

type pasteHandlerImpl struct {
	service services.PasteService
}

func NewPasteHandler(service services.PasteService) PasteHandler {
	return pasteHandlerImpl{service: service}
}

func (h pasteHandlerImpl) CreatePaste(c *gin.Context) {
	var paste dto.CreatePaste
	if err := c.ShouldBindJSON(&paste); err != nil {
		log.Printf("bindJSON error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	id, err := h.service.CreatePaste(c.Request.Context(), &paste)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"paste id": id})
}

func (h pasteHandlerImpl) GetPaste(c *gin.Context) {
	id := c.Param("id")
	paste, err := h.service.GetPasteByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
		return
	}

	if paste.Deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste has been deleted"})
		return
	}

	var input dto.GetPasteResponse
	_ = c.ShouldBindJSON(&input)
	input.IP = c.ClientIP()
	email, exist := c.Get("email")
	if exist {
		input.Email = email.(string)
	}

	err = h.service.GetPaste(&input, paste)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	pasteResp := dto.PasteResponse{
		ID:      paste.ID,
		Content: paste.Content,
	}

	c.JSON(http.StatusOK, gin.H{"paste": pasteResp})
}
