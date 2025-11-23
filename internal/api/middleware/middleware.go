package middleware

import (
	"insider-case/internal/constants"
	"insider-case/internal/pkg/logger"
	"insider-case/internal/pkg/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path

		c.Next()

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

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered", "error", err)
				response.InternalServerError(c, response.ErrorCodeInternalServerError, "Internal server error", nil)
				c.Abort()
			}
		}()
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, "+constants.HeaderAccessToken)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}

func JSONContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

func AuthMiddleware(accessToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		if accessToken == "" {
			logger.Warn("Access token is not configured, all requests will be rejected")
		}

		token := c.GetHeader(constants.HeaderAccessToken)
		if token == "" {
			token = c.GetHeader("x-access-token")
		}
		if token == "" {
			token = c.GetHeader("X-ACCESS-TOKEN")
		}
		if token == "" {
			token = c.Request.Header.Get(constants.HeaderAccessToken)
		}

		if token == "" {
			logger.Warn("Unauthorized access attempt - missing "+constants.HeaderAccessToken,
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"client_ip", c.ClientIP(),
				"user_agent", c.Request.UserAgent(),
				"status_code", http.StatusUnauthorized,
			)
			response.Unauthorized(c, response.ErrorCodeUnauthorizedMissingToken, constants.HeaderAccessToken+" header is required")
			c.Abort()
			return
		}

		if token != accessToken {
			logger.Warn("Unauthorized access attempt - invalid "+constants.HeaderAccessToken,
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"client_ip", c.ClientIP(),
				"user_agent", c.Request.UserAgent(),
				"status_code", http.StatusUnauthorized,
				"token_provided", token[:min(len(token), 10)]+"...",
			)
			response.Unauthorized(c, response.ErrorCodeUnauthorizedInvalidToken, "Invalid "+constants.HeaderAccessToken)
			c.Abort()
			return
		}

		c.Next()
	}
}
