package routes

import (
	_ "insider-case/docs" // Swagger documentation
	"insider-case/internal/api"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// setupSystemRoutes configures system-level routes (health, swagger)
func setupSystemRoutes(router *gin.Engine) {
	// Health check endpoint
	router.GET(api.HealthPath, healthCheck)

	// Swagger UI documentation
	router.GET(api.SwaggerPath, ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// healthCheck handles health check requests
func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
