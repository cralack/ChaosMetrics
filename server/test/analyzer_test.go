package test

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/service/analyzer"
	"github.com/cralack/ChaosMetrics/server/model/anres"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/utils"
)

func Test_analyzer(t *testing.T) {
	p := analyzer.NewAnalyzer(
		analyzer.WithLoc(riotmodel.NA1),
	)
	p.Analyze()
}

func Test_result_list(t *testing.T) {
	var (
		tmp      *anres.ChampionDetail
		aramList []*anres.ChampionDetail
		err      error
	)
	if err = db.Where("game_mode = ?", "CLASSIC").Where("version LIKE ?",
		"13.14%").Find(&aramList).Error; err != nil {
		t.Log("load aram analyzed result failed")
	} else {
		sort.Slice(aramList, func(i, j int) bool {
			return aramList[i].RankScore > aramList[j].RankScore
		})

		for i := 0; i < len(aramList); i++ {
			tmp = aramList[i]
			fmt.Printf("%d.%s:winrate:%0.2f score:%0.2f\n\r", i+1, tmp.Title, tmp.WinRate, tmp.RankScore)
		}
	}

}

func TestName(t *testing.T) {
	region := utils.ConvertLocationToRegionHost(riotmodel.TW2)
	puuid := "9qfUnTSxGzba5kG6Hk3t5SJtsO6D14AUpt5BMHhVyLuO17-FHsRpY_iJQOGmH9CED9DEtUX9QniDcw"
	startTime := time.Now().AddDate(-1, 0, 0).Unix() // one year ago unix
	endTime := time.Now().Unix()                     // cur time unix
	queryParams := fmt.Sprintf("startTime=%d&endTime=%d&start=0&count=%d",
		startTime, endTime, 20)
	url := fmt.Sprintf("%s/lol/match/v5/matches/by-puuid/%s/ids?%s",
		region, puuid, queryParams)
	buff, err := f.Get(url)
	if err != nil {
		t.Log(err)
	}
	t.Log(string(buff))
}

func Test_getCham(t *testing.T) {
	var (
		buff string
		err  error
		ctx  = context.Background()
	)
	buff = rdb.HGet(ctx, "/championlist", "1401").Val()
	keys := make([]string, 0)
	if err = json.Unmarshal([]byte(buff), &keys); err != nil {
		t.Log(err)
	}
	for i := range keys {
		keys[i] += "@1401"
	}

	values := rdb.HMGet(ctx, "/champions/zh_CN", keys...).Val()

	chams := make([]*riotmodel.ChampionDTO, 0, len(keys))
	for _, v := range values {
		var cham *riotmodel.ChampionDTO
		if err = json.Unmarshal([]byte(v.(string)), &cham); v == nil || err != nil {
			t.Log(err)
			continue
		}
		chams = append(chams, cham)
	}
	t.Log(len(chams))
}
