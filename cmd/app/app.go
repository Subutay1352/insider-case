package main

import (
	"context"
	"insider-case/internal/api/routes"
	"insider-case/internal/api/server"
	"insider-case/internal/config"
	"insider-case/internal/domain/message"
	"insider-case/internal/infrastructure/db"
	"insider-case/internal/infrastructure/httpclient"
	redisInfra "insider-case/internal/infrastructure/redis"
	"insider-case/internal/pkg/logger"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

type App struct {
	Config    *config.Config
	Service   *message.Service
	Scheduler *message.Scheduler
	Server    *http.Server
}

// NewApp initializes and returns the application
func NewApp(cfg *config.Config) (*App, error) {
	// Init DB
	database, err := db.InitDB(cfg)
	if err != nil {
		return nil, err
	}

	// Init Redis
	var cacheRepo message.CacheRepository
	var redisClient *redis.Client
	if rc, err := redisInfra.InitRedis(cfg); err == nil {
		redisClient = rc
		cacheRepo = redisInfra.NewCacheRepository(rc, cfg.Redis.TTL)
		logger.Info("Redis connection verified")
	} else {
		logger.Warn("Failed to initialize Redis, continuing without cache", "error", err)
	}

	// Init HTTP client
	webhookClient := httpclient.NewWebhookClient(cfg)

	// Init services
	messageRepo := db.NewRepository(database, cfg.Database.Type)
	messageService := message.NewService(
		messageRepo,
		cacheRepo,
		webhookClient,
		cfg.Scheduler.MessagesPerBatch,
		cfg.Message.MaxLength,
		cfg.Webhook.MaxRetryAttempts,
		cfg.Scheduler.RetryBaseDelay,
	)
	messageScheduler := message.NewScheduler(messageService, cfg.Scheduler.Interval, cfg.Scheduler.ProcessingTimeout)

	// Start scheduler if auto-start enabled
	if cfg.Scheduler.AutoStart {
		if err := messageScheduler.Start(); err != nil {
			logger.Warn("Failed to start scheduler", "error", err)
		}
	}

	// Setup routes and start server
	router := routes.SetupRoutes(messageService, messageScheduler, cfg, database, redisClient)
	srv := server.Start(router, &cfg.Server)

	return &App{
		Config:    cfg,
		Service:   messageService,
		Scheduler: messageScheduler,
		Server:    srv,
	}, nil
}

func (a *App) Shutdown() {
	logger.Info("Shutting down...")

	// Stop scheduler
	if a.Scheduler.IsRunning() {
		ctx, cancel := context.WithTimeout(context.Background(), a.Config.Scheduler.ShutdownTimeout)
		defer cancel()
		if err := a.Scheduler.StopAndWait(ctx); err != nil {
			logger.Error("Scheduler shutdown error", "error", err)
		}
	}

	// Shutdown server
	if err := server.Shutdown(a.Server, a.Config.Server.ShutdownTimeout); err != nil {
		logger.Error("Server shutdown error", "error", err)
	}

	time.Sleep(1 * time.Second)
}
