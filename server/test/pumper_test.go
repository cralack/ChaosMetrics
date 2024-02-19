package test

import (
	"testing"

	"github.com/cralack/ChaosMetrics/server/internal/service/pumper"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/utils"
)

func Test_pumper_update(t *testing.T) {
	p, _ := pumper.NewPumper(
		"1",
		pumper.WithLoc(riotmodel.NA1),
		pumper.WithEndMark(riotmodel.DIAMOND, 1),
		pumper.WithQues(riotmodel.RANKED_SOLO_5x5, riotmodel.RANKED_FLEX_SR),
	)
	p.StartEngine()
	p.UpdateAll()
}

func Test_pumper_fetch_match_byId(t *testing.T) {
	matchId := "TW2_99178853"
	region := utils.ConvertLocationToRegionHost(riotmodel.TW2)
	p, _ := pumper.NewPumper("1")
	p.FetchMatchByID(nil, region, matchId)
}

func Test_pumper_fetch_single_summoner(t *testing.T) {
	p, err := pumper.NewPumper("1",
		pumper.WithLoc(riotmodel.NA1),
	)
	if err != nil {
		t.Log(err)
	}
	go p.StartEngine()
	p.LoadSingleSummoner("Pink HairDryer", "na1")
}
