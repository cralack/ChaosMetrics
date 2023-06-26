package test

import (
	"strconv"
	"testing"

	"ChaosMetrics/server/config"
	"ChaosMetrics/server/global"
	"ChaosMetrics/server/pkg/db"
	"ChaosMetrics/server/pkg/vconf"
	"ChaosMetrics/server/pkg/zlog"
)

var err error

func goInit() {
	//set runtime envs
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

// basic pkg init test
func Test_config(t *testing.T) {
	goInit()
	conf := global.GVA_CONF
	global.GVA_LOG.Info(conf.Env)
	global.GVA_LOG.Info(conf.DirTree.WordDir)
	global.GVA_LOG.Info(conf.Dbconf.Driver)
	global.GVA_LOG.Info(strconv.Itoa(conf.LogConf.MaxSize))
	global.GVA_LOG.Info(conf.Riot.Apikey)
}
