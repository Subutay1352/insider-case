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

	// Initialize controllers
	senderController := controllers.NewSenderController(scheduler)
	messageController := controllers.NewMessageController(messageService, &cfg.Message)

	// System routes (no base path)
	setupSystemRoutes(router)

	// API v1 routes
	setupAPIRoutes(router, senderController, messageController, cfg)

	return router
}
