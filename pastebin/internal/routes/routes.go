package routes

import (
	"github.com/JuDyas/GolangPractice/pastebin/internal/auth"
	"github.com/JuDyas/GolangPractice/pastebin/internal/handlers"
	"github.com/JuDyas/GolangPractice/pastebin/internal/services"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, jwtSecret []byte, pasteHandler *handlers.PasteHandler, userService services.UserService) {
	v1 := r.Group("/v1")
	//Authorization
	v1.POST("/auth/register", handlers.Register(userService))
	v1.POST("/auth/login", handlers.Login(userService, jwtSecret))
	//Pasts
	v1.POST("/pastes", pasteHandler.CreatePaste)
	v1.GET("/pastes/:id", pasteHandler.GetPaste)
	//secured group
	authorize := v1.Group("/")
	authorize.Use(auth.AuthorizeMiddleware(jwtSecret))
	//authorize.DELETE("/pastes/:id", handlers.DeletePaste)
	//authorize.GET("/pastes", handlers.GetAllPastes)
}
