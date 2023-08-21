package test

import (
	"testing"
	
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/service/pumper"
	"github.com/cralack/ChaosMetrics/server/utils"
)

func Test_pumper_update(t *testing.T) {
	p := pumper.NewPumper(
		pumper.WithLoc(riotmodel.NA1),
		pumper.WithEndMark(riotmodel.DIAMOND, 1),
	)
	p.UpdateAll()
	
}

func Test_pumper_fetch_match_byId(t *testing.T) {
	matchId := "TW2_99178853"
	host := utils.ConvertPlatformToHost(riotmodel.TW2)
	p := pumper.NewPumper()
	p.FetchMatchByID(nil, host, matchId)
	
}