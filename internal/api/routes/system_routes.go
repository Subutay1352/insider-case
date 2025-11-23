package routes

import (
	"context"
	_ "insider-case/docs" // Swagger documentation
	"insider-case/internal/api"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// setupSystemRoutes configures system-level routes (health, swagger)
func setupSystemRoutes(router *gin.Engine, database *gorm.DB, redisClient *redis.Client) {
	// Health check endpoint (simple)
	router.GET(api.HealthPath, healthCheck)

	// System health check endpoint (DB + Redis)
	router.GET(api.HealthPathSystem, systemHealthCheck(database, redisClient))

	// Swagger UI documentation
	router.GET(api.SwaggerPath, ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// healthCheck handles simple health check requests
func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

// systemHealthCheck handles system health check with DB and Redis connection status
func systemHealthCheck(database *gorm.DB, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := gin.H{}

		// Check database connection
		if database != nil {
			sqlDB, err := database.DB()
			if err == nil {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				if err := sqlDB.PingContext(ctx); err == nil {
					status["db"] = "ok"
				} else {
					status["db"] = "error"
				}
			} else {
				status["db"] = "error"
			}
		} else {
			status["db"] = "not_configured"
		}

		// Check Redis connection
		if redisClient != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if err := redisClient.Ping(ctx).Err(); err == nil {
				status["redis"] = "ok"
			} else {
				status["redis"] = "error"
			}
		} else {
			status["redis"] = "not_configured"
		}

		c.JSON(200, status)
	}
}
