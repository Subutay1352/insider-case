package controllers

import (
	"insider-case/internal/config"
	"insider-case/internal/domain/message"
	"insider-case/internal/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MessageController handles message-related requests
type MessageController struct {
	service *message.Service
	config  *config.MessageConfig
}

// NewMessageController creates a new MessageController
func NewMessageController(service *message.Service, msgConfig *config.MessageConfig) *MessageController {
	return &MessageController{
		service: service,
		config:  msgConfig,
	}
}

// GetSentMessages retrieves a list of sent messages with pagination
func (c *MessageController) GetSentMessages(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", strconv.Itoa(c.config.DefaultLimit))
	offsetStr := ctx.DefaultQuery("offset", strconv.Itoa(c.config.DefaultOffset))

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = c.config.DefaultLimit
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = c.config.DefaultOffset
	}

	messages, total, err := c.service.GetSentMessages(ctx.Request.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(ctx, "Failed to retrieve sent messages", err)
		return
	}

	response.OK(ctx, "Sent messages retrieved successfully", gin.H{
		"messages": messages,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	})
}
