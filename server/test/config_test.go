package test

import (
	"strconv"
	"testing"

	"ChaosMetrics/server/global"
	_ "ChaosMetrics/server/init"
)

// basic pkg test
// include viper,zap,gorm
func Test_config(t *testing.T) {
	conf := global.GVA_CONF
	global.GVA_LOG.Info(conf.DirTree.WordDir)
	global.GVA_LOG.Info(conf.Dbconf.Driver)
	global.GVA_LOG.Info(strconv.Itoa(conf.LogConf.MaxSize))
	global.GVA_LOG.Info(conf.Riot.Apikey)
}
