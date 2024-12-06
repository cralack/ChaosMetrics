package test

import (
	"context"
	"testing"

	"github.com/cralack/ChaosMetrics/server/internal/service/pumper"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/utils"
)

func Test_pumper_update(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p, _ := pumper.NewPumper(
		"1",
		pumper.WithLoc(riotmodel.NA1),
		pumper.WithEndMark(riotmodel.DIAMOND, 1),
		pumper.WithQues(riotmodel.RANKED_SOLO_5x5, riotmodel.RANKED_FLEX_SR),
		pumper.WithMaxMatchCount(3),
		pumper.WithContext(ctx),
	)
	p.StartEngine()
	p.UpdateAll()
}

func Test_pumper_fetch_match_byId(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	matchId := "TW2_99178853"
	region := utils.ConvertLocationToRegionHost(riotmodel.TW2)
	sumnID := "PNqZUt2Ru1M_3x5aKyN3i73BK4uWHKT2owRa8qSBvyepeFU@RANKED_SOLO_5x5"
	p, _ := pumper.NewPumper("1", pumper.WithContext(ctx))
	p.FetchEntryBySumnID(sumnID, riotmodel.NA1)
	p.FetchMatchByID(nil, region, matchId)
}
