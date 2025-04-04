package middleware

import (
	"fmt"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start).Microseconds()
		if path != "/health" {
			global.ChaLogger.Info(path,
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
				zap.String("cost", fmt.Sprintf("%d ms", cost)),
			)
		} else {
			global.ChaLogger.Debug("health check")
		}
	}
}
