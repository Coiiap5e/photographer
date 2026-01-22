package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		if raw != "" {
			path = path + "?" + raw
		}

		latency := time.Since(start)

		logger.Info("request completed",
			"method", c.Request.Method,
			"path", path,
			"ip", c.ClientIP(),
			"user-agent", c.Request.UserAgent(),
			"status", c.Writer.Status(),
			"latency", latency,
		)
	}
}
