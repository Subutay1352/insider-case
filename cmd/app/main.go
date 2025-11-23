package main

import (
	"os"
	"os/signal"
	"syscall"

	"insider-case/internal/config"
	"insider-case/internal/pkg/logger"
)

func main() {
	cfg := config.Load()
	logger.Init(cfg.Env)

	app, err := NewApp(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize application", "error", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.Shutdown()
}
