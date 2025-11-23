package routes

import (
	"insider-case/internal/api"
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
	v1 := router.Group(api.APIV1BasePath)
	// Apply authentication middleware to all API routes
	v1.Use(middleware.AuthMiddleware(cfg.AccessToken))
	{
		// Sender endpoints
		sender := v1.Group(api.SenderBasePath)
		{
			sender.POST(api.StartSchedulerPath, senderController.Start)
			sender.POST(api.StopSchedulerPath, senderController.Stop)
			sender.GET(api.StatusSchedulerPath, senderController.Status)
		}

		// Message endpoints
		messages := v1.Group(api.MessagesBasePath)
		{
			messages.GET(api.SentMessagesPath, messageController.GetSentMessages)
		}
	}
}
