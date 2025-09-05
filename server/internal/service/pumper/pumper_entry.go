package pumper

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/cralack/ChaosMetrics/server/utils/scheduler"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type entryTask struct {
	Tier  string
	Rank  string
	Queue string
}

func (p *Pumper) UpdateEntries() {
	for _, loc := range p.stgy.Loc {
		p.loadEntrie(loc)
		for _, que := range p.stgy.Que {
			// Generate URLs
			go p.createEntriesURL(loc, que)
		}
	}
	// blocking until handle final
	<-p.Exit
}

func (p *Pumper) loadEntrie(location riotmodel.LOCATION) {
	loc, _ := utils.ConvertLocationToLoHoSTR(location)
	ctx := context.Background()
	key := fmt.Sprintf("/entry/%s", loc)
	// init if local map doesn't exist
	if _, has := p.entryMap[loc]; !has {
		p.entryMap[loc] = make(map[string]*riotmodel.LeagueEntryDTO)
	}

	// load if redis cache exist
	redisMap := make(map[string]*riotmodel.LeagueEntryDTO)
	if size := p.rdb.HLen(ctx, key).Val(); size != 0 {
		kvmap := p.rdb.HGetAll(ctx, key).Val()
		for k, v := range kvmap {
			var tmp riotmodel.LeagueEntryDTO
			if err := json.Unmarshal([]byte(v), &tmp); err != nil {
				p.logger.Error("load entry form redis cache failed", zap.Error(err))
			} else {
				redisMap[k] = &tmp
			}
		}
	}

	// load if gorm db data exist
	var entries []*riotmodel.LeagueEntryDTO
	if err := p.db.Where("loc = ?", loc).Find(&entries).Error; err != nil {
		p.logger.Error("load entry from gorm db failed", zap.Error(err))
	}
	if len(entries) != 0 {
		for _, e := range entries {
			k := fmt.Sprintf("%s@%s", e.SummonerID, e.QueType)
			// assign to localmap
			if _, has := p.entryMap[loc][k]; !has {
				p.entryMap[loc][k] = e
			}
		} // assign to redis
	}

	// check local && redis map diff
	tmp := make([]*riotmodel.LeagueEntryDTO, 0, p.stgy.MaxSize)
	for k, e := range redisMap {
		if _, has := p.entryMap[loc][k]; !has {
			p.entryMap[loc][k] = e
			tmp = append(tmp, e)
		}
	}
	var mx uint = 0
	for k, e := range p.entryMap[loc] {
		if _, has := redisMap[k]; !has {
			tmp = append(tmp, e)
		}
		if mx < e.ID {
			mx = e.ID
		}
	}

	loCode := uint(utils.ConvertLocStrToLocation(loc))
	p.lock.Lock()
	p.entrieIdx[loCode] += (mx+1)%(loCode*1e9) + (loCode * 1e9)
	p.lock.Unlock()
	// store diff to db
	p.handleEntries(tmp, loc)
}

func (p *Pumper) createEntriesURL(loc riotmodel.LOCATION, que riotmodel.QUECODE) {
	var (
		url  string
		tier riotmodel.TIER
		rank uint
	)
	locStr, host := utils.ConvertLocationToLoHoSTR(loc)
	queStr := utils.ConvertQueToQueSTR(que)
	// p.loadEntrie(locStr)

	// generate BEST URL task
	for tier = riotmodel.CHALLENGER; tier <= riotmodel.MASTER; tier++ {
		if tier > p.stgy.TestEndMark1 || (tier == p.stgy.TestEndMark1 && rank > p.stgy.TestEndMark2) {
			return
		}
		t, r := ConvertRankToStr(tier, 1)
		url = fmt.Sprintf("%s/lol/league/v4/%sleagues/by-queue/%s",
			host, strings.ToLower(t), queStr)
		p.scheduler.Push(&scheduler.Task{
			Type: bestEntryTypeKey,
			Loc:  locStr,
			URL:  url,
			Data: &entryTask{
				Tier:  t,
				Rank:  r,
				Queue: queStr,
			},
		})
	}
	// generate MORTAL URL task
	for tier = riotmodel.DIAMOND; tier <= riotmodel.IRON; tier++ {
		for rank = 1; rank <= 4; rank++ {
			if tier > p.stgy.TestEndMark1 || (tier == p.stgy.TestEndMark1 && rank > p.stgy.TestEndMark2) {
				return
			}
			t, r := ConvertRankToStr(tier, rank)
			url = fmt.Sprintf("%s/lol/league/v4/entries/%s/%s/%s",
				host, queStr, t, r)
			p.scheduler.Push(&scheduler.Task{
				Type: mortalEntryTypeKey,
				Loc:  locStr,
				URL:  url,
				Data: &entryTask{
					Tier:  t,
					Rank:  r,
					Queue: queStr,
				},
			})
		}
	}
}

func (p *Pumper) handleEntries(entries []*riotmodel.LeagueEntryDTO, loc string) {
	if len(entries) == 0 {
		return
	}
	tmp := make([]*riotmodel.LeagueEntryDTO, 0, p.stgy.MaxSize)
	loCode := utils.ConvertLocStrToLocation(loc)

	p.lock.Lock()
	defer p.lock.Unlock()

	for _, entry := range entries {
		// entry doenst exist:create new data
		k := fmt.Sprintf("%s@%s", entry.SummonerID, entry.QueType)
		if e, has := p.entryMap[loc][k]; !has {
			p.entryMap[loc][k] = entry
			entry.ID = p.entrieIdx[loCode]
			p.entrieIdx[loCode]++
		} else {
			// update old data(will wipe redis's gorm date)
			entry.ID = e.ID
		}
		tmp = append(tmp, entry)
	}

	if len(tmp) == 0 {
		return
	}
	// check oversize && split
	var chunks [][]*riotmodel.LeagueEntryDTO
	if len(tmp) > p.stgy.MaxSize {
		totalSize := len(tmp)
		chunkSize := p.stgy.MaxSize
		for i := 0; i < totalSize; i += chunkSize {
			end := i + chunkSize
			if end > totalSize {
				end = totalSize
			}
			chunks = append(chunks, tmp[i:end])
		}
	} else {
		chunks = append(chunks, tmp)
	}

	// send to DB handler
	for _, chunk := range chunks {
		p.out <- &DBResult{
			Type:  EntryTypeKey,
			Brief: chunk[0].Tier + " " + chunk[0].Rank,
			Data:  chunk,
		}
	}
}

func (p *Pumper) cacheEntries(entries []*riotmodel.LeagueEntryDTO, loc string) {
	ctx := context.Background()
	key := fmt.Sprintf("/entry/%s", loc)
	pipe := p.rdb.Pipeline()
	pipe.Expire(ctx, key, p.stgy.LifeTime)
	cmds := make([]*redis.IntCmd, 0, len(entries))
	for _, e := range entries {
		field := fmt.Sprintf("%s@%s", e.SummonerID, e.QueType)
		cmds = append(cmds, pipe.HSet(ctx, key, field, e))
	}
	if _, err := pipe.Exec(ctx); err != nil {
		p.logger.Error("redis store entry failed", zap.Error(err))
	}
}

func (p *Pumper) FetchEntryBySumnID(sumnID string, loc riotmodel.LOCATION) {
	var (
		url    string
		locStr string
		host   string
	)
	locStr, host = utils.ConvertLocationToLoHoSTR(loc)

	url = fmt.Sprintf("%s/lol/league/v4/entries/by-summoner/%s", host, sumnID)
	p.scheduler.Push(&scheduler.Task{
		Type:     mortalEntryTypeKey,
		Priority: true,
		Loc:      locStr,
		URL:      url,
		Data: &entryTask{
			Tier:  "",
			Rank:  "",
			Queue: "",
		},
	})
}
