package routes

import (
	_ "insider-case/docs" // Swagger documentation

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// setupSystemRoutes configures system-level routes (health, swagger)
func setupSystemRoutes(router *gin.Engine) {
	// Health check endpoint
	router.GET("/health", healthCheck)

	// Swagger UI documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// healthCheck handles health check requests
func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
