package http

import (
	"auth-server/internal/utils"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, handlers *AuthHandlers, jwtService *utils.JWTService) {
	// Публичные маршруты
	authGroup := router.Group("/")
	{
		authGroup.POST("/register", handlers.Register)
		authGroup.POST("/login", handlers.Login)
		authGroup.POST("/refresh", handlers.Refresh)
	}

	// Защищённые маршруты
	protected := router.Group("/")
	protected.Use(AuthMiddleware(jwtService))
	{
		protected.GET("/protected", handlers.ProtectedExample)
	}
}
