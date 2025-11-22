package middleware

import (
	"insider-case/internal/pkg/logger"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path

		// Process request
		c.Next()

		// Calculate latency
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		logger.Info("HTTP Request",
			"method", method,
			"path", path,
			"status", statusCode,
			"latency_ms", duration.Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}

// ErrorHandlingMiddleware handles panics and errors
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered", "error", err)
				c.JSON(500, gin.H{"error": "Internal server error"})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// CORSMiddleware adds CORS headers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Access-Token")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}

// JSONContentTypeMiddleware sets JSON content type
func JSONContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

// AuthMiddleware validates X-Access-Token header
func AuthMiddleware(accessToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow OPTIONS requests for CORS preflight
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Warn if access token is not configured
		if accessToken == "" {
			logger.Warn("Access token is not configured, all requests will be rejected")
		}

		// Get token from header (case-insensitive)
		// HTTP headers are case-insensitive, but we check multiple variations for compatibility
		token := c.GetHeader("X-Access-Token")
		if token == "" {
			token = c.GetHeader("x-access-token")
		}
		if token == "" {
			token = c.GetHeader("X-ACCESS-TOKEN")
		}
		// Also check directly from request headers (Request.Header.Get is case-insensitive)
		if token == "" {
			token = c.Request.Header.Get("X-Access-Token")
		}

		// If no token provided or token doesn't match, return 401
		if token == "" {
			logger.Warn("Unauthorized access attempt - missing X-Access-Token",
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"client_ip", c.ClientIP(),
				"user_agent", c.Request.UserAgent(),
				"status_code", http.StatusUnauthorized,
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: X-Access-Token header is required",
			})
			c.Abort()
			return
		}

		// Check if token matches
		if token != accessToken {
			logger.Warn("Unauthorized access attempt - invalid X-Access-Token",
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"client_ip", c.ClientIP(),
				"user_agent", c.Request.UserAgent(),
				"status_code", http.StatusUnauthorized,
				"token_provided", token[:min(len(token), 10)]+"...", // Log first 10 chars for security
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: Invalid X-Access-Token",
			})
			c.Abort()
			return
		}

		// Token is valid, continue
		c.Next()
	}
}
