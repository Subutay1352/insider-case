package routes

import (
	"context"
	_ "insider-case/docs" // Swagger documentation
	"insider-case/internal/constants"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// setupSystemRoutes configures system-level routes (health, swagger)
func setupSystemRoutes(router *gin.Engine, database *gorm.DB, redisClient *redis.Client) {
	// Health check endpoint (DB + Redis)
	router.GET(constants.HealthPath, healthCheck(database, redisClient))

	// Swagger UI documentation
	router.GET(constants.SwaggerPath, ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// healthCheck handles health check with DB and Redis connection status
func healthCheck(database *gorm.DB, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := gin.H{}

		// Check database connection
		status["db"] = checkDatabaseHealth(database)

		// Check Redis connection
		status["redis"] = checkRedisHealth(redisClient)

		c.JSON(200, status)
	}
}

// checkDatabaseHealth checks database connection health
func checkDatabaseHealth(database *gorm.DB) string {
	if database == nil {
		return "not_configured"
	}

	sqlDB, err := database.DB()
	if err != nil {
		return "error"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return "error"
	}

	return "ok"
}

// checkRedisHealth checks Redis connection health
func checkRedisHealth(redisClient *redis.Client) string {
	if redisClient == nil {
		return "not_configured"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return "error"
	}

	return "ok"
}
