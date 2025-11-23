package routes

import (
	"insider-case/internal/api/controllers"
	"insider-case/internal/api/middleware"
	"insider-case/internal/config"

	"github.com/gin-gonic/gin"
)

// setupAPIRoutes configures API v1 routes
func setupAPIRoutes(
	router *gin.Engine,
	senderController *controllers.SenderController,
	messageController *controllers.MessageController,
	cfg *config.Config,
) {
	v1 := router.Group("/api/v1")
	// Apply authentication middleware to all API routes
	v1.Use(middleware.AuthMiddleware(cfg.AccessToken))
	{
		// Sender endpoints
		sender := v1.Group("/sender")
		{
			sender.POST("/startScheduler", senderController.Start)
			sender.POST("/stopScheduler", senderController.Stop)
			sender.GET("/statusScheduler", senderController.Status)
		}

		// Message endpoints
		messages := v1.Group("/messages")
		{
			messages.GET("/sent", messageController.GetSentMessages)
		}
	}
}
