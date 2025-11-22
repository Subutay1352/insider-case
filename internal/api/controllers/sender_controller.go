package controllers

import (
	"insider-case/internal/domain/message"
	"insider-case/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// SenderController handles sender-related requests
type SenderController struct {
	scheduler *message.Scheduler
}

// NewSenderController creates a new SenderController
func NewSenderController(scheduler *message.Scheduler) *SenderController {
	return &SenderController{
		scheduler: scheduler,
	}
}

// Start starts the automatic message sending scheduler
// @Summary      Start scheduler
// @Description  Starts the automatic message sending scheduler
// @Tags         sender
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Scheduler started successfully"
// @Failure      400  {object}  map[string]interface{}  "Scheduler is already running"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /api/v1/sender/start [post]
func (c *SenderController) Start(ctx *gin.Context) {
	if c.scheduler.IsRunning() {
		response.BadRequest(ctx, "Scheduler is already running")
		return
	}

	if err := c.scheduler.Start(); err != nil {
		response.InternalServerError(ctx, "Failed to start scheduler", err)
		return
	}

	response.OK(ctx, "Scheduler started successfully", nil)
}

// Stop stops the automatic message sending scheduler
// @Summary      Stop scheduler
// @Description  Stops the automatic message sending scheduler
// @Tags         sender
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Scheduler stopped successfully"
// @Failure      400  {object}  map[string]interface{}  "Scheduler is not running"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /api/v1/sender/stop [post]
func (c *SenderController) Stop(ctx *gin.Context) {
	if !c.scheduler.IsRunning() {
		response.BadRequest(ctx, "Scheduler is not running")
		return
	}

	if err := c.scheduler.Stop(); err != nil {
		response.InternalServerError(ctx, "Failed to stop scheduler", err)
		return
	}

	response.OK(ctx, "Scheduler stopped successfully", nil)
}

// Status returns the current status of the scheduler
// @Summary      Get scheduler status
// @Description  Returns the current status of the scheduler
// @Tags         sender
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Scheduler status with is_running field"
// @Router       /api/v1/sender/status [get]
func (c *SenderController) Status(ctx *gin.Context) {
	response.OK(ctx, "Scheduler status retrieved", gin.H{
		"is_running": c.scheduler.IsRunning(),
	})
}
