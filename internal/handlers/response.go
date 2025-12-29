package handlers

import (
	"net/http"
	"todo-go-backend/internal/errors"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse represents a standardized success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// handleError handles application errors and returns appropriate HTTP response
func handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		c.JSON(appErr.StatusCode, ErrorResponse{
			Error:   appErr.Error(),
			Message: appErr.Message,
		})
		return
	}

	// Generic error
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error:   "Internal server error",
		Message: err.Error(),
	})
}

// handleSuccess returns a standardized success response
func handleSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	response := SuccessResponse{
		Message: message,
	}
	if data != nil {
		response.Data = data
	}
	c.JSON(statusCode, response)
}

// handleValidationError handles Gin validation errors
func handleValidationError(c *gin.Context, err error) {
	handleError(c, errors.NewInvalidInputError(err.Error()))
}

