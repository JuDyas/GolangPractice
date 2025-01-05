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
	GetPaste(c *gin.Context)
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

func (h pasteHandlerImpl) GetPaste(c *gin.Context) {
	pasteVal, exists := c.Get("paste")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "paste not found (ctx)"})
		return
	}

	paste, ok := pasteVal.(*models.Paste)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid paste object"})
		return
	}

	if paste.Deleted {
		c.JSON(http.StatusNotFound, gin.H{"error": "paste has been deleted"})
		return
	}

	var input models.InputPaste
	_ = c.ShouldBindJSON(&input)
	input.IP = c.ClientIP()
	email, exist := c.Get("email")
	if exist {
		input.Email = email.(string)
	}

	err := h.service.GetPaste(c.Request.Context(), &input, paste)
	//TODO: Разобраться с ошибками
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pasteDTL := models.PasteDTl{
		ID:      paste.ID,
		Content: paste.Content,
	}

	c.JSON(http.StatusOK, gin.H{"paste": pasteDTL})
}
