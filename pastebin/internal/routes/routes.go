package routes

import (
	"github.com/JuDyas/GolangPractice/pastebin_new/internal/auth"
	"github.com/JuDyas/GolangPractice/pastebin_new/internal/handlers"
	"github.com/JuDyas/GolangPractice/pastebin_new/internal/services"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userService services.UserService, pasteService services.PasteService, jwtSecret []byte, pasteHandler handlers.PasteHandler, adminHandler handlers.AdminPasteHandler) {
	v1 := r.Group("/v1")
	v1.Use(auth.AuthoriseMiddleware(jwtSecret))
	//Authorisation
	v1.POST("/auth/register", handlers.Register(userService))
	v1.POST("/auth/login", handlers.Login(userService, jwtSecret))
	//Pastes
	v1.POST("/pastes", pasteHandler.CreatePaste)
	v1.GET("/:id", pasteHandler.GetPaste)
	//Admin group
	admin := v1.Group("/admin")
	admin.DELETE("/pastes/:id", adminHandler.DeletePaste)
	admin.POST("/pastes", adminHandler.ListPastes)
}
