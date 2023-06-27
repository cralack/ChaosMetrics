package test

import (
	"strconv"
	"testing"

	"ChaosMetrics/server/global"
)

// basic pkg init test
func Test_config(t *testing.T) {
	conf := global.GVA_CONF
	global.GVA_LOG.Info(conf.Env)
	global.GVA_LOG.Info(conf.DirTree.WorkDir)
	global.GVA_LOG.Info(conf.Dbconf.Driver)
	global.GVA_LOG.Info(strconv.Itoa(conf.LogConf.MaxSize))
	global.GVA_LOG.Info(strconv.Itoa(int(conf.Fetcher.Timeout)))
}
