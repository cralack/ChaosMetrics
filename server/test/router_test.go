package test

import (
	"fmt"
	"testing"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/internal/service/router"
	"go.uber.org/zap"
)

func Test_router(t *testing.T) {
	routeConf := global.ChaConf.Router
	logger.Debug(fmt.Sprintf("prifix:%s", routeConf.RouterPrefix))

	r := router.New()

	if err := r.Run(":8080"); err != nil {
		logger.Error("router failed", zap.Error(err))
	}
}
