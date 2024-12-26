package routes

import (
	"github.com/JuDyas/GolangPractice/pastebin/internal/auth"
	"github.com/JuDyas/GolangPractice/pastebin/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, jwtSecret []byte, port string) {
	v1 := r.Group("/v1")
	//Authorization
	v1.POST("/auth/register", handlers.Register(jwtSecret))
	v1.POST("/auth/login", handlers.Login(jwtSecret))
	//Pasts
	v1.POST("/pastes", handlers.CreatePaste)
	v1.GET("/pastes/:id", handlers.GetPaste)
	//secured group
	authorize := v1.Group("/")
	authorize.Use(auth.AuthorizeMiddleware(jwtSecret))
	authorize.DELETE("/pastes/:id", handlers.DeletePaste)
	authorize.GET("/pastes", handlers.GetAllPastes)
}
