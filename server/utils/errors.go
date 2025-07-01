package utils

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"visit-tracker-api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// HandleError handles different types of errors and responds appropriately
func HandleError(c *gin.Context, err error, message string) {
	requestID := c.GetString("request_id")
	
	fields := logrus.Fields{
		"request_id": requestID,
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"ip":         c.ClientIP(),
	}

	// Check error type and respond accordingly
	switch {
	case errors.Is(err, sql.ErrNoRows):
		apiErr := &middleware.APIError{
			Code:       "NOT_FOUND",
			Message:    "Resource not found",
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
		LogError(err, message, fields)
		c.Error(apiErr)
		return

	case isValidationError(err):
		apiErr := &middleware.APIError{
			Code:       "VALIDATION_ERROR",
			Message:    "Validation failed",
			Details:    err.Error(),
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
		LogWarn(message, fields)
		c.Error(apiErr)
		return

	case isDatabaseError(err):
		apiErr := &middleware.APIError{
			Code:       "DATABASE_ERROR",
			Message:    "Database operation failed",
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
		LogError(err, message, fields)
		c.Error(apiErr)
		return

	default:
		apiErr := &middleware.APIError{
			Code:       "INTERNAL_SERVER_ERROR",
			Message:    "Internal server error",
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
		LogError(err, message, fields)
		c.Error(apiErr)
		return
	}
}

// HandleValidationError handles validation errors specifically
func HandleValidationError(c *gin.Context, err error, field string) {
	requestID := c.GetString("request_id")
	
	apiErr := &middleware.APIError{
		Code:    "VALIDATION_ERROR",
		Message: "Validation failed",
		Details: map[string]string{
			"field": field,
			"error": err.Error(),
		},
		StatusCode: http.StatusBadRequest,
		Err:        err,
	}

	LogWarn("Validation error", logrus.Fields{
		"request_id": requestID,
		"field":      field,
		"error":      err.Error(),
	})

	c.Error(apiErr)
}

// HandleDatabaseError handles database errors specifically
func HandleDatabaseError(c *gin.Context, err error, operation string) {
	requestID := c.GetString("request_id")
	
	fields := logrus.Fields{
		"request_id": requestID,
		"operation":  operation,
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
	}

	if errors.Is(err, sql.ErrNoRows) {
		apiErr := &middleware.APIError{
			Code:       "NOT_FOUND",
			Message:    "Resource not found",
			StatusCode: http.StatusNotFound,
			Err:        err,
		}
		LogWarn("Resource not found", fields)
		c.Error(apiErr)
		return
	}

	apiErr := &middleware.APIError{
		Code:       "DATABASE_ERROR",
		Message:    "Database operation failed",
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
	LogError(err, "Database error", fields)
	c.Error(apiErr)
}

// Helper functions to identify error types
func isValidationError(err error) bool {
	// Add logic to identify validation errors
	// This could check for specific error types or error messages
	return false // Placeholder
}

func isDatabaseError(err error) bool {
	// Add logic to identify database errors
	// This could check for specific database error types
	return true // For now, assume most errors are database-related
}

// Success response helper
func JSONSuccess(c *gin.Context, data interface{}) {
	requestID := c.GetString("request_id")
	
	response := gin.H{
		"data":       data,
		"request_id": requestID,
		"timestamp":  getCurrentTimestamp(),
	}
	
	c.JSON(http.StatusOK, response)
}

// Created response helper
func JSONCreated(c *gin.Context, data interface{}) {
	requestID := c.GetString("request_id")
	
	response := gin.H{
		"data":       data,
		"request_id": requestID,
		"timestamp":  getCurrentTimestamp(),
	}
	
	c.JSON(http.StatusCreated, response)
}

func getCurrentTimestamp() string {
	return time.Now().Format("2006-01-02T15:04:05Z07:00")
} 