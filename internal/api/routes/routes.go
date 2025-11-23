package routes

import (
	"insider-case/internal/api/controllers"
	"insider-case/internal/config"
	"insider-case/internal/domain/message"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	messageService *message.Service,
	scheduler *message.Scheduler,
	cfg *config.Config,
	database *gorm.DB,
	redisClient *redis.Client,
) *gin.Engine {
	router := gin.Default()

	// Initialize controllers
	senderController := controllers.NewSenderController(scheduler)
	messageController := controllers.NewMessageController(messageService, &cfg.Message)

	// System routes (no base path)
	setupSystemRoutes(router, database, redisClient)

	// API v1 routes
	setupAPIRoutes(router, senderController, messageController, cfg)

	return router
}
