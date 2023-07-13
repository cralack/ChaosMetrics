package pumper

import (
	"context"
	"encoding/json"
	"fmt"
	
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"go.uber.org/zap"
)

func (p *Pumper) UpdateSumonerMatch() {
	exit := make(chan struct{})
	for _, loc := range p.stgy.Loc {
		// Generate URLs
		go p.createSummonerWork(loc)
	}
	<-exit
}

func (p *Pumper) createSummonerWork(loc uint) {
	// var (
	// 	path string
	// 	tier uint
	// 	div  uint
	// )
	// locStr, host := utils.ConvertHostURL(loc)
	// for tier = riotmodel.CHALLENGER; tier <= riotmodel.MASTER; tier++ {
	// 	t, r := ConvertRankToStr(tier, 1)
	//
	// }
}

func (p *Pumper) loadEntrie(loc string) {
	ctx := context.Background()
	res := make(map[string]*riotmodel.LeagueEntryDTO)
	// load if local map exist
	if entries, has := p.entriesDic[loc]; !has || len(entries) == 0 {
		key := fmt.Sprintf("/entry/%s", loc)
		
		// load if redis cache exist
		if size := p.rdb.HLen(ctx, key).Val(); size != 0 {
			kvmap := p.rdb.HGetAll(ctx, key).Val()
			for k, v := range kvmap {
				var tmp riotmodel.LeagueEntryDTO
				if err := json.Unmarshal([]byte(v), &tmp); err != nil {
					p.logger.Error("load entry form redis cache failed", zap.Error(err))
				} else {
					res[k] = &tmp
				}
			}
			
			// update local?
		}
		
		// load if gorm db data exist
	}
	
	return
}
