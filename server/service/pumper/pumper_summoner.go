package pumper

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"time"
	
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/service/scheduler"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

type summonerTask struct {
	summonerID string
	url        string
	summoner   *riotmodel.SummonerDTO
}

func (p *Pumper) UpdateSumoner(exit chan struct{}) {
	for _, loc := range p.stgy.Loc {
		// Generate URLs
		go p.createSummonerURL(loc)
	}
	go p.fetchSummoner()
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
	var max uint = 0
	for k, s := range p.sumnMap[loc] {
		if _, has := redisMap[k]; !has {
			tmp = append(tmp, s)
		}
		if max < s.ID {
			max = s.ID
		}
	}
	loCode := utils.ConverHostLoCode(loc)
	p.lock.Lock()
	p.summonerIdx[loCode] += (max+1)%(loCode*1e9) + (loCode * 1e9)
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
			Key: "summoner",
			Loc: loc,
			Data: &summonerTask{
				summonerID: sID,
				url:        url,
			},
		})
	}
	// expand from self
	for sID, sumn := range p.sumnMap[loc] {
		if sumn.Name == "" {
			url := fmt.Sprintf("%s/lol/summoner/v4/summoners/%s", host, sID)
			p.scheduler.Push(&scheduler.Task{
				Key: "summoner",
				Loc: loc,
				Data: &summonerTask{
					summonerID: sID,
					url:        url,
				},
			})
		}
	}
	// finish signal
	p.scheduler.Push(&scheduler.Task{
		Key:  "summoner",
		Loc:  loc,
		Data: nil,
	})
	return
}

func (p *Pumper) fetchSummoner() {
	var (
		buff []byte
		err  error
	)
	
	defer func() {
		if err := recover(); err != nil {
			p.logger.Panic("fetcher panic",
				zap.Any("err", err),
				zap.String("stack", string(debug.Stack())))
		}
	}()
	
	for {
		// <-p.Rater
		req := p.scheduler.Pull()
		switch req.Key {
		case "summoner":
			if req.Data == nil {
				// send finish signal
				p.out <- &ParseResult{
					Type:  "finish",
					Brief: "summoner",
					Data:  nil,
				}
				// *need release scheduler resource*
				return
			} else {
				data := req.Data.(*summonerTask)
				if buff, err = p.fetcher.Get(data.url); err != nil || buff == nil {
					p.logger.Error(fmt.Sprintf("fetch summonerID %s failed", data.summonerID), zap.Error(err))
					// fetch again
					if req.Retry < p.stgy.Retry {
						req.Retry++
						p.scheduler.Push(req)
					}
					continue
				}
				var sumn *riotmodel.SummonerDTO
				if err = json.Unmarshal(buff, &sumn); err != nil {
					p.logger.Error(fmt.Sprintf("unmarshal json to %s failed",
						"sumnDTO"), zap.Error(err))
				}
				p.handleSummoner(req.Loc, sumn)
			}
		}
	}
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
		p.out <- &ParseResult{
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
			p.out <- &ParseResult{
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
		delta int
		rate  float32
	)
	ticker := time.NewTicker(time.Second * 30)
	total = len(p.entryMap[loc])
	for {
		<-ticker.C
		delta = len(p.sumnMap[loc]) - cur
		cur = len(p.sumnMap[loc])
		rate = float32(cur) / float32(total)
		if delta > 0 {
			p.logger.Info(fmt.Sprintf("fetch %s %05.02f%% (%04d/%04d) summoners",
				loc, rate*100, cur, total))
		}
	}
}