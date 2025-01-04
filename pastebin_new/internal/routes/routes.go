package routes

import (
	"github.com/JuDyas/GolangPractice/pastebin_new/internal/auth"
	"github.com/JuDyas/GolangPractice/pastebin_new/internal/handlers"
	"github.com/JuDyas/GolangPractice/pastebin_new/internal/services"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userService services.UserService, jwtSecret []byte, pasteHandler handlers.PasteHandler) {
	v1 := r.Group("/v1")
	//Authorisation
	v1.POST("/auth/register", handlers.Register(userService))
	v1.POST("/auth/login", handlers.Login(userService, jwtSecret))
	//Pastes
	v1.POST("/pastes", pasteHandler.CreatePaste)
	//secured group
	authorize := v1.Group("/")
	authorize.Use(auth.AuthoriseMiddleware(jwtSecret))
}
