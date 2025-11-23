package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, statusCode int, result *SuccessResult) {
	c.JSON(statusCode, result)
}

func ErrorResponse(c *gin.Context, statusCode int, result *ErrorResult) {
	c.JSON(statusCode, result)
}

func BadRequest(c *gin.Context, code ErrorCode, message string) {
	ErrorResponse(c, http.StatusBadRequest, &ErrorResult{
		Code:    code,
		Message: message,
	})
}

func Unauthorized(c *gin.Context, code ErrorCode, message string) {
	ErrorResponse(c, http.StatusUnauthorized, &ErrorResult{
		Code:    code,
		Message: message,
	})
}

func InternalServerError(c *gin.Context, code ErrorCode, message string, err error) {
	result := &ErrorResult{
		Code:    code,
		Message: message,
	}
	if err != nil {
		result.Details = err.Error()
	}
	ErrorResponse(c, http.StatusInternalServerError, result)
}

func OK(c *gin.Context, code SuccessCode, message string, data interface{}) {
	SuccessResponse(c, http.StatusOK, &SuccessResult{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func Created(c *gin.Context, code SuccessCode, message string, data interface{}) {
	SuccessResponse(c, http.StatusCreated, &SuccessResult{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
