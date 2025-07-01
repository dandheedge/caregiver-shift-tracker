package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// responseWriter is a wrapper around gin.ResponseWriter to capture response body
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggingMiddleware returns a gin middleware for logging requests and responses
func LoggingMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Capture request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewReader(requestBody))
		}

		// Wrap response writer to capture response body
		responseBodyWriter := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer([]byte{}),
		}
		c.Writer = responseBodyWriter

		// Process request
		c.Next()

		// Log request and response
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		logEntry := logger.WithFields(logrus.Fields{
			"method":        c.Request.Method,
			"path":          c.Request.URL.Path,
			"query":         c.Request.URL.RawQuery,
			"status":        statusCode,
			"duration":      duration,
			"duration_ms":   duration.Milliseconds(),
			"ip":           c.ClientIP(),
			"user_agent":   c.Request.UserAgent(),
			"request_id":   c.GetString("request_id"),
		})

		// Add request body for non-GET requests (be careful with sensitive data)
		if c.Request.Method != "GET" && len(requestBody) > 0 && len(requestBody) < 1024 {
			logEntry = logEntry.WithField("request_body", string(requestBody))
		}

		// Add response body for errors (status >= 400)
		if statusCode >= 400 && responseBodyWriter.body.Len() > 0 && responseBodyWriter.body.Len() < 1024 {
			logEntry = logEntry.WithField("response_body", responseBodyWriter.body.String())
		}

		// Log with appropriate level based on status code
		switch {
		case statusCode >= 500:
			logEntry.Error("Internal server error")
		case statusCode >= 400:
			logEntry.Warn("Client error")
		case statusCode >= 300:
			logEntry.Info("Redirection")
		default:
			logEntry.Info("Request completed")
		}
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := generateRequestID()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// Simple request ID generator
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
} 