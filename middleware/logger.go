package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Logger(c *gin.Context) {
	logger := log.Logger.With().Str("request_id", c.Query("request_id")).Logger()
	c.Set("logger", logger)

	start := time.Now()
	c.Next()
	latency := time.Since(start)

	method := c.Request.Method
	path := c.Request.URL.Path
	logger.Info().
		Str("method", method).
		Str("path", path).
		Int("status", c.Writer.Status()).
		Str("client_ip", c.ClientIP()).
		Dur("latency", latency).
		Int("bytes", c.Writer.Size()).
		Send()
}
