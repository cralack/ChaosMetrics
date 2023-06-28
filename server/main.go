package main

import (
	"github.com/cralack/ChaosMetrics/server/global"
	_ "github.com/cralack/ChaosMetrics/server/init"
	"os"
)

func main() {
	wd, _ := os.Getwd()
	global.GVA_LOG.Info(wd)
}
