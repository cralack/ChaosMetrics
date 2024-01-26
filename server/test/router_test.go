package test

import (
	"testing"

	"github.com/cralack/ChaosMetrics/server/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Test_router(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.GinLogger())

	r.GET("/ping", func(context *gin.Context) {
		context.JSONP(200, gin.H{
			"message": "pong",
		})
	})
	if err := r.Run(":8080"); err != nil {
		logger.Error("router failed", zap.Error(err))
	}
}
