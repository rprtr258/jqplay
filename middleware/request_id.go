package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID(c *gin.Context) {
	requestID := c.Request.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = uuid.New().String()
	}
	c.Set("request_id", requestID)
	c.Next()
}
