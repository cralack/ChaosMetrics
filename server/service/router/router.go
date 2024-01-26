package router

import (
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	switch global.ChaEnv {

	}

	router := gin.New()
	router.Use(
		middleware.GinLogger(),
		gin.Recovery(),
	)

	gin.Default()
	return router
}
