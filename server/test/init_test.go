package test

import (
	"strconv"
	"testing"

	"github.com/cralack/ChaosMetrics/server/global"
)

// basic pkg init test
func Test_config(t *testing.T) {
	conf := global.GVA_CONF
	logger.Info(conf.Env)
	logger.Info(conf.DirTree.WorkDir)
	logger.Info(conf.Dbconf.Driver)
	logger.Info(strconv.Itoa(conf.LogConf.MaxSize))
	logger.Info(strconv.Itoa(int(conf.Fetcher.Timeout)))
}
