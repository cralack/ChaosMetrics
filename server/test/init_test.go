package test

import (
	"strconv"
	"testing"

	"github.com/cralack/ChaosMetrics/server/global"
)

// basic pkg init test
func Test_config(t *testing.T) {
	conf := global.GVA_CONF
	logger.Debug(conf.Env)
	logger.Debug(conf.DirTree.WorkDir)
	logger.Debug(conf.Dbconf.Driver)
	logger.Debug(strconv.Itoa(conf.LogConf.MaxSize))
	logger.Debug(strconv.Itoa(int(conf.Fetcher.Timeout)))
}
