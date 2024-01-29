package router

import (
	"github.com/cralack/ChaosMetrics/server/app/api"
	"github.com/cralack/ChaosMetrics/server/app/middleware"
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	switch global.ChaEnv {

	}

	router := gin.New()
	router.Use(
		middleware.GinLogger(),
		gin.Recovery(),
	)

	api.RegiserRoutes(router)

	return router
}
