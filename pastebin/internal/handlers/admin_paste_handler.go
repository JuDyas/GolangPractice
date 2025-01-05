package handlers

import (
	"net/http"

	"github.com/JuDyas/GolangPractice/pastebin_new/internal/services"
	"github.com/gin-gonic/gin"
)

type AdminPasteHandler interface {
	DeletePaste(c *gin.Context)
	ListPastes(c *gin.Context)
}

type adminHandlerImpl struct {
	service services.AdminPasteService
}

func NewAdminHandler(service services.AdminPasteService) AdminPasteHandler {
	return &adminHandlerImpl{service: service}
}

func (h *adminHandlerImpl) DeletePaste(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.SoftDeletePaste(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "paste delete successfully"})
}

func (h *adminHandlerImpl) ListPastes(c *gin.Context) {
	var input struct {
		Content   string                 `json:"content,omitempty"`
		CreatedAt map[string]interface{} `json:"created_at,omitempty"`
		SortBy    string                 `json:"sortBy,omitempty"`
		Page      int                    `json:"page,omitempty"`
		Limit     int                    `json:"limit,omitempty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if input.SortBy == "" {
		input.SortBy = "created_at"
	}
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 10
	}

	filters := make(map[string]interface{})
	if input.Content != "" {
		filters["content"] = input.Content
	}
	if input.CreatedAt != nil {
		filters["created_at"] = input.CreatedAt
	}

	pastes, err := h.service.GetAllPastes(c.Request.Context(), filters, input.SortBy, input.Page, input.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pastes": pastes})
}
