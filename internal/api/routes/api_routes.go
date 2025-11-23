package routes

import (
	"insider-case/internal/api/controllers"
	"insider-case/internal/api/middleware"
	"insider-case/internal/config"
	"insider-case/internal/constants"

	"github.com/gin-gonic/gin"
)

// setupAPIRoutes configures API v1 routes
func setupAPIRoutes(
	router *gin.Engine,
	senderController *controllers.SenderController,
	messageController *controllers.MessageController,
	cfg *config.Config,
) {
	v1 := router.Group(constants.APIV1BasePath)
	// Apply authentication middleware to all API routes
	v1.Use(middleware.AuthMiddleware(cfg.AccessToken))
	{
		// Sender endpoints
		sender := v1.Group(constants.SenderBasePath)
		{
			sender.POST(constants.StartSchedulerPath, senderController.Start)
			sender.POST(constants.StopSchedulerPath, senderController.Stop)
			sender.GET(constants.StatusSchedulerPath, senderController.Status)
		}

		// Message endpoints
		messages := v1.Group(constants.MessagesBasePath)
		{
			messages.GET(constants.SentMessagesPath, messageController.GetSentMessages)
		}
	}
}
