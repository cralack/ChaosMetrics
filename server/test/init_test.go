package test

import (
	"strconv"
	"testing"

	"github.com/cralack/ChaosMetrics/server/internal/global"
)

// basic pkg init test
func Test_config(t *testing.T) {
	conf := global.GvaConf
	logger.Debug(conf.Env)
	logger.Debug(conf.DirTree.WorkDir)
	logger.Debug(conf.Dbconf.Driver)
	logger.Debug(strconv.Itoa(conf.LogConf.MaxSize))
	logger.Debug(strconv.Itoa(int(conf.Fetcher.Timeout)))
}
