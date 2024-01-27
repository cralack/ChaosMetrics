package test

import (
	"testing"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/internal/service/router"
	"go.uber.org/zap"
)

func Test_router(t *testing.T) {
	prefix := global.ChaConf.System.RouterPrefix
	logger.Debug(prefix)

	r := router.New()

	if err := r.Run(":8080"); err != nil {
		logger.Error("router failed", zap.Error(err))
	}
}
