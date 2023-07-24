package main

import (
	"os"
	"os/signal"
	
	"github.com/cralack/ChaosMetrics/server/global"
	_ "github.com/cralack/ChaosMetrics/server/init"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/service/pumper"
	"go.uber.org/zap"
)

func main() {
	p := pumper.NewPumper(
		"",
		pumper.WithEndMark(riotmodel.DIAMOND, 4),
		pumper.WithLoc(riotmodel.TW2),
	)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	go p.UpdateAll()
	sig := <-quit
	global.GVA_LOG.Info("catch signal,exiting...",
		zap.Any("signal", sig))
}
