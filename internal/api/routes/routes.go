package routes

import (
	"insider-case/internal/api/controllers"
	"insider-case/internal/config"
	"insider-case/internal/domain/message"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	messageService *message.Service,
	scheduler *message.Scheduler,
	cfg *config.Config,
) *gin.Engine {
	router := gin.Default()

	// Controllers
	senderController := controllers.NewSenderController(scheduler)
	messageController := controllers.NewMessageController(messageService, &cfg.Message)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Sender routes
	router.POST("/sender/start", senderController.Start)
	router.POST("/sender/stop", senderController.Stop)
	router.GET("/sender/status", senderController.Status)

	// Message routes
	router.GET("/messages/sent", messageController.GetSentMessages)

	return router
}
