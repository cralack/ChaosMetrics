package test

import (
	"testing"
	
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/service/updater"
)

func Test_updater(t *testing.T) {
	updtr := updater.NewRiotUpdater()
	puuid := "F4fFtqehQLBj8U5sKBZF--k-7akbtb1IX790lRd4whPI4pXDAuVyfswHetg2lz_kMe2NJ0gUo5EIig"
	matches, err := updtr.UpdateSummonerMatch(riotmodel.TW2, puuid, 20)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		t.Log(matches[i].Metadata.MetaMatchID)
		t.Log(matches[i].Info.GameVersion)
	}
}

func Test_update_best(t *testing.T) {
	updtr := updater.NewRiotUpdater()
	var loc, que, tier uint
	loc = riotmodel.TW2
	que = riotmodel.RANKED_SOLO_5x5
	for tier = riotmodel.CHALLENGER; tier <= riotmodel.MASTER; tier++ {
		res, err := updtr.UpdateBetserLeague(loc, tier, que)
		if err != nil {
			t.Fatal(err)
		} else {
			t.Log("update succeed ", res[0].Tier, len(res))
		}
	}
}

func Test_update_all(t *testing.T) {
	updtr := updater.NewRiotUpdater()
	var loc, que, tier, div, page uint
	loc = riotmodel.TW2
	que = riotmodel.RANKED_SOLO_5x5
	// batch query
	// for tier = riotmodel.DIAMOND; tier <= riotmodel.IRON; tier++ {
	// 	for div = 1; div <= 4; div++ {
	// 		for page = 1; page <= 1; page++ {
	// 			res, err := updtr.UpdateMortalLeague(loc, tier, div, que, page)
	// 			if err != nil {
	// 				t.Fatal(err)
	// 			} else {
	// 				t.Log("update succeed ", res[0].Tier, len(res))
	// 			}
	// 		}
	// 	}
	// }
	// single query
	tier = riotmodel.DIAMOND
	div = 1
	page = 1
	res, err := updtr.UpdateMortalLeague(loc, tier, div, que, page)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("update succeed ", len(res))
	}
}
