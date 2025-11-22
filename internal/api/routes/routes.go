package routes

import (
	_ "insider-case/docs" // Swagger documentation
	"insider-case/internal/api/controllers"
	"insider-case/internal/api/middleware"
	"insider-case/internal/config"
	"insider-case/internal/domain/message"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

// setupSystemRoutes configures system-level routes (health, swagger)
func setupSystemRoutes(router *gin.Engine) {
	// Health check endpoint
	router.GET("/health", healthCheck)

	// Swagger UI documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

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
			sender.POST("/start", senderController.Start)
			sender.POST("/stop", senderController.Stop)
			sender.GET("/status", senderController.Status)
		}

		// Message endpoints
		messages := v1.Group("/messages")
		{
			messages.GET("/sent", messageController.GetSentMessages)
		}
	}
}

// healthCheck handles health check requests
func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
