package init

import (
	"ChaosMetrics/server/config"
	"ChaosMetrics/server/global"
	"ChaosMetrics/server/pkg/db"
	"ChaosMetrics/server/pkg/vconf"
	"ChaosMetrics/server/pkg/zlog"
)

var err error

func init() {
	//set runtime env
	global.GVA_ENV = global.TEST_ENV
	global.GVA_CONF = config.New()

	//setup config service
	global.GVA_VP, err = vconf.Viper()
	if err != nil {
		panic(err)
	}

	//setup logger service
	global.GVA_LOG, err = zlog.Zap(global.GVA_ENV)
	if err != nil {
		panic(err)
	}

	//setup orm service
	global.GVA_DB, err = db.GetDB()
	if err != nil {
		panic(err)
	}

	global.GVA_LOG.Info("env pkg init succeed")
}
