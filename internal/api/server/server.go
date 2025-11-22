package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"insider-case/internal/api/middleware"
	"insider-case/internal/config"
	"insider-case/internal/pkg/logger"

	"github.com/gin-gonic/gin"
)

// Start starts the HTTP server in a goroutine
func Start(router *gin.Engine, cfg *config.ServerConfig) *http.Server {
	// Apply middleware
	router.Use(middleware.ErrorHandlingMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.JSONContentTypeMiddleware())
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		logger.Info("Server starting", "address", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	return srv
}

// Shutdown gracefully shuts down the server
func Shutdown(srv *http.Server, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return srv.Shutdown(ctx)
}
