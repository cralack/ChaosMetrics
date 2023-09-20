package main

import (
	"github.com/cralack/ChaosMetrics/server/cmd"
	_ "github.com/cralack/ChaosMetrics/server/init"
)

func main() {
	// p := pumper.NewPumper(
	// 	pumper.WithEndMark(riotmodel.MASTER, 1),
	// 	pumper.WithLoc(riotmodel.NA1),
	// )
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt, os.Kill)

	if err := cmd.RunCommand(); err != nil {
		panic(err)
	}

	// sig := <-quit
	// global.GVA_LOG.Info("catch signal,exiting...",
	// 	zap.Any("signal", sig))
}
