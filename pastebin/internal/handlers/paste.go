package handlers

import (
	"net/http"
	"time"

	"github.com/JuDyas/GolangPractice/pastebin/internal/db"

	"github.com/JuDyas/GolangPractice/pastebin/models"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: Add zap logger for errors
func CreatePaste(c *gin.Context) {
	var (
		paste      models.Paste
		collection = db.Client.Database("pastebin").Collection("pastes")
	)
	if err := c.ShouldBindJSON(&paste); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paste.ID = primitive.NewObjectID().Hex()
	paste.CreatedAt = time.Now()
	_, err := collection.InsertOne(c.Request.Context(), paste)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	//TODO: Поменять
	c.JSON(http.StatusCreated, gin.H{"id": paste.ID})
}

func GetPaste(c *gin.Context) {
	var (
		id         = c.Param("id")
		paste      models.Paste
		collection = db.Client.Database("pastebin").Collection("pastes")
	)
	err := collection.FindOne(c.Request.Context(), bson.M{"_id": id}).Decode(&paste)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//TODO: Изменить проверку для бесконечных паст
	if !paste.ExpiresAt.IsZero() && time.Now().After(paste.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"error": "paste expired"})
		return
	}
	c.JSON(http.StatusOK, paste)
}
