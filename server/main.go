package main

import (
	"os"

	"github.com/cralack/ChaosMetrics/server/global"
	_ "github.com/cralack/ChaosMetrics/server/init"
)

func main() {
	wd, _ := os.Getwd()
	global.GVA_LOG.Info(wd)
}
