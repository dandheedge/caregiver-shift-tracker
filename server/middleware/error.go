package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error     ErrorDetail `json:"error"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp string      `json:"timestamp"`
}

// ErrorDetail contains error information
type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// APIError represents an application error
type APIError struct {
	Code       string
	Message    string
	Details    interface{}
	StatusCode int
	Err        error
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// NewAPIError creates a new API error
func NewAPIError(code, message string, statusCode int, details interface{}) *APIError {
	return &APIError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Details:    details,
	}
}

// Common API errors
var (
	ErrInternalServer = NewAPIError("INTERNAL_SERVER_ERROR", "Internal server error", http.StatusInternalServerError, nil)
	ErrBadRequest     = NewAPIError("BAD_REQUEST", "Bad request", http.StatusBadRequest, nil)
	ErrNotFound       = NewAPIError("NOT_FOUND", "Resource not found", http.StatusNotFound, nil)
	ErrUnauthorized   = NewAPIError("UNAUTHORIZED", "Unauthorized", http.StatusUnauthorized, nil)
	ErrForbidden      = NewAPIError("FORBIDDEN", "Forbidden", http.StatusForbidden, nil)
	ErrValidation     = NewAPIError("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, nil)
)

// ErrorHandlerMiddleware handles panics and errors
func ErrorHandlerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := c.GetString("request_id")
				
				// Log the panic with stack trace
				logger.WithFields(logrus.Fields{
					"request_id": requestID,
					"method":     c.Request.Method,
					"path":       c.Request.URL.Path,
					"ip":         c.ClientIP(),
					"panic":      err,
					"stack":      string(debug.Stack()),
				}).Error("Panic recovered")

				// Return standardized error response
				errorResponse := ErrorResponse{
					Error: ErrorDetail{
						Code:    "INTERNAL_SERVER_ERROR",
						Message: "Internal server error",
					},
					RequestID: requestID,
					Timestamp: getCurrentTimestamp(),
				}

				c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse)
			}
		}()

		c.Next()

		// Handle errors set by handlers
		if len(c.Errors) > 0 {
			requestID := c.GetString("request_id")
			err := c.Errors.Last()

			// Log the error
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"ip":         c.ClientIP(),
				"error":      err.Error(),
			}).Error("Request error")

			// Check if it's an APIError
			if apiErr, ok := err.Err.(*APIError); ok {
				errorResponse := ErrorResponse{
					Error: ErrorDetail{
						Code:    apiErr.Code,
						Message: apiErr.Message,
						Details: apiErr.Details,
					},
					RequestID: requestID,
					Timestamp: getCurrentTimestamp(),
				}
				c.AbortWithStatusJSON(apiErr.StatusCode, errorResponse)
				return
			}

			// Generic error response
			errorResponse := ErrorResponse{
				Error: ErrorDetail{
					Code:    "INTERNAL_SERVER_ERROR",
					Message: "Internal server error",
				},
				RequestID: requestID,
				Timestamp: getCurrentTimestamp(),
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse)
		}
	}
}

// Helper function to get current timestamp
func getCurrentTimestamp() string {
	return time.Now().Format("2006-01-02T15:04:05Z07:00")
} 