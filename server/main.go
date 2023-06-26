package main

import (
	"ChaosMetrics/server/global"
	_ "ChaosMetrics/server/init"
	"os"
)

func main() {
	wd, _ := os.Getwd()
	global.GVA_LOG.Info(wd)

}
