package test

import (
	"testing"

	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/service/pumper"
	"github.com/cralack/ChaosMetrics/server/utils"
)

func Test_pumper_update(t *testing.T) {
	p, _ := pumper.NewPumper(
		"1",
		pumper.WithLoc(riotmodel.NA1),
		pumper.WithEndMark(riotmodel.EMERALD, 1),
	)
	exit := make(chan struct{})
	p.StartEngine(exit)
	p.UpdateAll(exit)
	<-exit
}

func Test_pumper_fetch_match_byId(t *testing.T) {
	matchId := "TW2_99178853"
	region := utils.ConvertLocToRegion(riotmodel.TW2)
	p, _ := pumper.NewPumper("1")
	p.FetchMatchByID(nil, region, matchId)
}

func Test_pumper_fetch_single_summoner(t *testing.T) {
	p, err := pumper.NewPumper("1")
	if err != nil {
		t.Log(err)
	}
	// exit := make(chan struct{})
	// go p.StartEngine(exit)
	// p.loadSummoners("tw2")
	p.LoadSingleSummoner("Mes", "tw2")
}
