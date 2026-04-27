package http

import (
	"auth-server/internal/delivery/http/middleware"
	"auth-server/internal/utils"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authHandlers *AuthHandlers, jwtService *utils.JWTService) {
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandlers.Register)
			auth.POST("/login", authHandlers.Login)
			auth.GET("/profile", middleware.AuthMiddleware(jwtService), profileHandler) // profileHandler нужно определить или создать
		}
	}
}

// временный profileHandler
func profileHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(200, gin.H{"user_id": userID})
}
