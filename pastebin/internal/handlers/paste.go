package handlers

import (
	"net/http"

	"github.com/JuDyas/GolangPractice/pastebin/internal/servises"

	"github.com/JuDyas/GolangPractice/pastebin/models"
	"github.com/gin-gonic/gin"
)

func CreatePasteHandler(pasteService *servises.PasteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var paste models.Paste
		if err := c.ShouldBindJSON(&paste); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		if err := pasteService.CreatePaste(&paste); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": paste.ID})
	}
}

//func GetPaste(c *gin.Context) {
//	var (
//		id         = c.Param("id")
//		paste      models.Paste
//		collection = db.Client.Database("pastebin").Collection("pastes")
//	)
//	err := collection.FindOne(c.Request.Context(), bson.M{"_id": id}).Decode(&paste)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	if paste.TTL > 0 && paste.CreatedAt.Unix()+int64(paste.TTL) <= time.Now().Unix() {
//		c.JSON(http.StatusGone, gin.H{"error": "paste expired"})
//		return
//	}
//
//	pass, ok := c.GetQuery("pass")
//	if paste.Password != "" && pass != paste.Password {
//		if !ok {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "password is required"})
//		} else {
//			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
//		}
//		return
//	}
//
//	publicPaste := models.PublicPaste{
//		ID:   paste.ID,
//		Text: paste.Text,
//	}
//
//	c.JSON(http.StatusOK, publicPaste)
//}
