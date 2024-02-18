package router

import (
	"github.com/cralack/ChaosMetrics/server/app/api"
	"github.com/cralack/ChaosMetrics/server/app/middleware"
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var addr string

func init() {
	Cmd.Flags().StringVar(&addr, "addr", ":8080", "set router listen addr")
}

var Cmd = &cobra.Command{
	Use:   "router",
	Short: "launches a router to handle frontend requests",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		r := New()
		return r.Run(addr)
	},
}

func New() *gin.Engine {
	switch global.ChaEnv {
	case global.TestEnv:
		gin.SetMode(gin.TestMode)
	case global.DevEnv:
		gin.SetMode(gin.DebugMode)
	case global.ProductEnv:
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(
		middleware.GinLogger(),
		gin.Recovery(),
	)

	api.RegisterRoutes(router)

	return router
}
