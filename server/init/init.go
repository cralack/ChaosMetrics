package init

import (
	"github.com/cralack/ChaosMetrics/server/config"
	"github.com/cralack/ChaosMetrics/server/global"
	"github.com/cralack/ChaosMetrics/server/pkg/db"
	"github.com/cralack/ChaosMetrics/server/pkg/vconf"
	"github.com/cralack/ChaosMetrics/server/pkg/zlog"
)

var err error

func init() {
	// set runtime envs
	global.GVA_CONF = config.New()

	// setup config service
	global.GVA_VP, err = vconf.Viper()
	if err != nil {
		panic(err)
	}

	// setup logger service
	global.GVA_LOG, err = zlog.Zap(global.GVA_ENV)
	if err != nil {
		panic(err)
	}

	// setup orm service
	global.GVA_DB, err = db.GetDB()
	if err != nil {
		panic(err)
	}

	global.GVA_LOG.Info("env pkg init succeed")
}
