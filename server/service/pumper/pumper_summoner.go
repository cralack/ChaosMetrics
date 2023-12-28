package pumper

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/cralack/ChaosMetrics/server/utils/scheduler"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

type summonerTask struct {
	summonerID string
	summoner   *riotmodel.SummonerDTO
}

func (p *Pumper) UpdateSumoner(exit chan struct{}) {
	for _, loc := range p.stgy.Loc {
		// Generate URLs
		go p.createSummonerURL(loc)
	}
	// go p.fetchSummoner()
	<-exit
}

func (p *Pumper) loadSummoner(loc string) {
	ctx := context.Background()
	key := fmt.Sprintf("/summoner/%s", loc)
	// init if local map doesn't exist
	if _, has := p.sumnMap[loc]; !has {
		p.sumnMap[loc] = make(map[string]*riotmodel.SummonerDTO)
	}

	// load if redis cache exist
	redisMap := make(map[string]*riotmodel.SummonerDTO)
	if size := p.rdb.HLen(ctx, key).Val(); size != 0 {
		kvmap := p.rdb.HGetAll(ctx, key).Val()
		for k, v := range kvmap {
			var tmp riotmodel.SummonerDTO
			if err := json.Unmarshal([]byte(v), &tmp); err != nil {
				p.logger.Error("load entry form redis cache failed", zap.Error(err))
			} else {
				redisMap[k] = &tmp
			}
		}
	}
	// load if gorm db data exist
	var summoners []*riotmodel.SummonerDTO
	if err := p.db.Where("loc = ?", loc).Find(&summoners).Preload(clause.Associations).Error; err != nil {
		p.logger.Error("load summoner from gorm db failed", zap.Error(err))
	}
	if len(summoners) != 0 {
		for _, s := range summoners {
			// assign to localmap
			if _, has := p.sumnMap[loc][s.MetaSummonerID]; !has {
				p.sumnMap[loc][s.MetaSummonerID] = s
			}
		} // assign to redis
		p.cacheSummoners(summoners, loc)
	}
	// check local && redis map diff
	tmp := make([]*riotmodel.SummonerDTO, 0, p.stgy.MaxSize)
	for k, s := range redisMap {
		if _, has := p.sumnMap[loc][k]; !has {
			p.sumnMap[loc][k] = s
			tmp = append(tmp, s)
		}
	}
	var mx uint = 0
	for k, s := range p.sumnMap[loc] {
		if _, has := redisMap[k]; !has {
			tmp = append(tmp, s)
		}
		if mx < s.ID {
			mx = s.ID
		}
	}
	loCode := utils.ConverHostLoCode(loc)
	p.lock.Lock()
	p.summonerIdx[loCode] += (mx+1)%(loCode*1e9) + (loCode * 1e9)
	p.lock.Unlock()
	p.handleSummoner(loc, tmp...)
}

func (p *Pumper) createSummonerURL(loCode uint) {
	loc, host := utils.ConvertHostURL(loCode)
	p.loadSummoner(loc)
	go p.summonerCounter(loc)
	// expand from entry
	for sID, entry := range p.entryMap[loc] {
		curTier, curRank := ConvertStrToRank(entry.Tier, entry.Rank)
		if _, has := p.sumnMap[loc][sID]; has ||
			(curTier > p.stgy.TestEndMark[0] || (curTier == p.stgy.TestEndMark[0] && curRank > p.stgy.TestEndMark[1])) {
			continue
		}
		url := fmt.Sprintf("%s/lol/summoner/v4/summoners/%s", host, sID)
		p.scheduler.Push(&scheduler.Task{
			Type: "summoner",
			Loc:  loc,
			URL:  url,
			Data: &summonerTask{
				summonerID: sID,
			},
		})
	}
	// expand from self
	for sID, sumn := range p.sumnMap[loc] {
		if sumn.Name == "" {
			url := fmt.Sprintf("%s/lol/summoner/v4/summoners/%s", host, sID)
			p.scheduler.Push(&scheduler.Task{
				Type: "summoner",
				Loc:  loc,
				URL:  url,
				Data: &summonerTask{
					summonerID: sID,
				},
			})
		}
	}
	// finish signal
	p.scheduler.Push(&scheduler.Task{
		Type: "summoner",
		Loc:  loc,
		Data: nil,
	})
	return
}

func (p *Pumper) handleSummoner(loc string, summoners ...*riotmodel.SummonerDTO) {
	if len(summoners) == 0 {
		return
	}
	// set up redis pipe
	// tmp := make([]*riotmodel.SummonerDTO, 0, p.stgy.MaxSize)
	loCode := utils.ConverHostLoCode(loc)
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, sumn := range summoners {
		if s, has := p.sumnMap[loc][sumn.MetaSummonerID]; !has || s.ID < 1e9*loCode {
			sumn.ID = p.summonerIdx[loCode]
			p.summonerIdx[loCode]++
		} else {
			sumn.ID = s.ID
		}
		p.sumnMap[loc][sumn.MetaSummonerID] = sumn
		sumn.Loc = loc
		p.rdb.HSet(context.Background(), "/summoner/"+loc, sumn.MetaSummonerID, sumn)
	}

	// check oversize && split
	if len(summoners) < p.stgy.MaxSize {
		p.out <- &DBResult{
			Type: "summoners",
			Data: summoners,
		}
	} else {
		totalSize := len(summoners)
		chunkSize := p.stgy.MaxSize
		for i := 0; i < totalSize; i += chunkSize {
			end := i + chunkSize
			if end > totalSize {
				end = totalSize
			}
			p.out <- &DBResult{
				Type: "summoners",
				Data: summoners[i:end],
			}
		}
	}

	return
}

func (p *Pumper) cacheSummoners(summoners []*riotmodel.SummonerDTO, loc string) {
	ctx := context.Background()
	key := fmt.Sprintf("/summoner/%s", loc)
	pipe := p.rdb.Pipeline()
	pipe.Expire(ctx, key, p.stgy.LifeTime)
	cmds := make([]*redis.IntCmd, 0, len(summoners))
	for _, s := range summoners {
		cmds = append(cmds, pipe.HSet(ctx, key, s.MetaSummonerID, s))
	}
	if _, err := pipe.Exec(ctx); err != nil {
		p.logger.Error("redis store summoner failed", zap.Error(err))
	}
}

func (p *Pumper) summonerCounter(loc string) {
	var (
		total int
		cur   int
		pre   int
		rate  float32
	)
	ticker := time.NewTicker(time.Second * 15)
	checkTicker := time.NewTicker(time.Millisecond * 100)
	total = len(p.entryMap[loc])
	for {
		select {
		case <-checkTicker.C:
			cur = len(p.sumnMap[loc])
			if cur == total {
				p.logger.Info("all ready fetch 100% summoners")
				return
			}
		case <-ticker.C:
			rate = float32(cur) / float32(total)
			if cur-pre > 0 {
				p.logger.Info(fmt.Sprintf("fetch %s %05.02f%% (%04d/%04d) summoners",
					loc, rate*100, cur, total))
			}
			pre = len(p.sumnMap[loc])
		}
	}
}
