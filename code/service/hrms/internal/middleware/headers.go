package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"hrms/pkg/errors"
)

const (
	tentantIDHeader = "X-Tenant-ID"
	clientIDHeader = "X-Client-ID"
)

// Headers is a middleware that validates required and optional headers
func Headers(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip for health check endpoint
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		// Check for required X-Tenant-ID header
		tentantID := c.GetHeader(tentantIDHeader)
		if tentantID == "" {
			err := errors.New("MISSING_HEADER", "X-Tenant-ID header is required")
			logger.WithError(err).Error("Missing required header")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Get optional X-Client-ID header
		clientID := c.GetHeader(clientIDHeader)

		// Set headers in context for later use
		c.Set("tenantID", tentantID)
		if clientID != "" {
			c.Set("clientID", clientID)
		}

		c.Next()
	}
}
