package pumper

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
	
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/service/scheduler"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type entryTask struct {
	Tier string
	Rank string
	URL  string
}

func (p *Pumper) UpdateEntries(exit chan struct{}) {
	for _, loc := range p.stgy.Loc {
		for _, que := range p.stgy.Que {
			// Generate URLs
			go p.CreateEntriesURL(loc, que)
		}
	}
	// go Watch_ETCD
	go p.FetchEntry()
	// blocking until handle final
	// go p.HandleResult(exit)
	<-exit
}

func (p *Pumper) CreateEntriesURL(loc, que uint) {
	var (
		path string
		tier uint
		rank uint
	)
	locStr, host := utils.ConvertHostURL(loc)
	queStr := getQueueString(que)
	p.loadEntrie(locStr)
	
	// generate BEST URL task
	for tier = riotmodel.CHALLENGER; tier <= riotmodel.MASTER; tier++ {
		t, r := ConvertRankToStr(tier, 1)
		
		path = fmt.Sprintf("/lol/league/v4/%sleagues/by-queue/%s",
			strings.ToLower(t), queStr)
		p.scheduler.Push(&scheduler.Task{
			Key: "bestEntry",
			Loc: locStr,
			Data: &entryTask{
				Tier: t,
				Rank: r,
				URL:  host + path,
			},
		})
	}
	// generate MORTAL URL task
	for tier = riotmodel.DIAMOND; tier <= riotmodel.IRON; tier++ {
		for rank = 1; rank <= 4; rank++ {
			t, r := ConvertRankToStr(tier, rank)
			path = fmt.Sprintf("/lol/league/v4/entries/%s/%s/%s", queStr, t, r)
			if tier > p.stgy.EndMark[0] || (tier == p.stgy.EndMark[0] && rank > p.stgy.EndMark[1]) {
				break
			}
			p.scheduler.Push(&scheduler.Task{
				Key: "mortalEntry",
				Loc: locStr,
				Data: &entryTask{
					Tier: t,
					Rank: r,
					URL:  host + path,
				},
			})
		}
	}
}

func (p *Pumper) FetchEntry() {
	var (
		page    int
		buff    []byte
		err     error
		list    *riotmodel.LeagueListDTO
		entries []*riotmodel.LeagueEntryDTO
		endTier string
		endRank string
	)
	endTier, endRank = ConvertRankToStr(p.stgy.EndMark[0], p.stgy.EndMark[1])
	// catch panic
	defer func() {
		if err := recover(); err != nil {
			p.logger.Error("fetcher panic",
				zap.Any("err", err),
				zap.String("stack", string(debug.Stack())))
		}
	}()
	
	for {
		req := p.scheduler.Pull()
		// fetch and parse
		switch req.Key {
		case "bestEntry":
			// rate limiter
			<-p.rater
			t := req.Data.(*entryTask)
			// api: /lol/league/v4/{BEST}leagues/by-queue/{queue}
			if buff, err = p.fetcher.Get(t.URL); err != nil {
				p.logger.Error(fmt.Sprintf("fetch %s %s failed", t.Tier, t.Rank),
					zap.Error(err))
				// fetch again
				p.scheduler.Push(req)
				continue
			}
			if err = json.Unmarshal(buff, &list); err != nil {
				p.logger.Error(fmt.Sprintf("unmarshal json to %s failed",
					"LeagueListDTO"), zap.Error(err))
			}
			entries = list.Entries
			for _, e := range entries {
				e.Tier = t.Tier
				e.Loc = req.Loc
			}
			list = nil
			p.logger.Info(fmt.Sprintf("all %d %s data fetch done",
				len(entries), t.Tier))
			p.handleEntries(entries, req.Loc)
			p.cacheEntries(entries, req.Loc)
		
		case "mortalEntry":
			page = 0
			for {
				page++
				// rate limiter
				<-p.rater
				t := req.Data.(*entryTask)
				// api: /lol/league/v4/entries/{queue}/{tier}/{division}
				if buff, err = p.fetcher.Get(fmt.Sprintf("%s?page=%s",
					t.URL, strconv.Itoa(page))); err != nil {
					p.logger.Error(fmt.Sprintf("fetch %s %s failed", t.Tier, t.Rank),
						zap.Error(err))
				}
				if err = json.Unmarshal(buff, &entries); err != nil {
					p.logger.Error(fmt.Sprintf("unmarshal json to %s failed",
						"LeagueEntryDTO"), zap.Error(err))
				} else {
					p.logger.Info(fmt.Sprintf("fetch %s %s page %d done", t.Tier, t.Rank, page))
				}
				for _, e := range entries {
					e.Loc = req.Loc
				}
				// send finish signal
				if len(entries) == 0 {
					p.logger.Info(fmt.Sprintf("all %s %s data fetch done at page %d",
						t.Tier, t.Rank, page))
					if t.Tier == endTier && t.Rank == endRank {
						p.out <- &ParseResult{
							Type:  "finish",
							Brief: "entry",
							Data:  nil,
						}
						// *need release scheduler resource*
						return
					}
					
					break
				}
				p.handleEntries(entries, req.Loc)
				p.cacheEntries(entries, req.Loc)
			}
			
		}
	}
}

func (p *Pumper) handleEntries(entries []*riotmodel.LeagueEntryDTO, loc string) {
	
	tmp := make([]*riotmodel.LeagueEntryDTO, 0, p.stgy.MaxSize)
	loCode := utils.ConverHostLoCode(loc)
	
	p.lock.Lock()
	for _, entry := range entries {
		// entry doenst exist:create new data
		if e, has := p.entryMap[loc][entry.SummonerID]; !has {
			p.entryMap[loc][entry.SummonerID] = entry
			entry.ID = p.entrieIdx[loCode]
			p.entrieIdx[loCode]++
		} else {
			// update old data
			entry.ID = e.ID
		}
		tmp = append(tmp, entry)
	}
	p.lock.Unlock()
	
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
		p.out <- &ParseResult{
			Type:  "entry",
			Brief: chunk[0].Tier + " " + chunk[0].Rank,
			Data:  chunk,
		}
	}
	return
}

func (p *Pumper) cacheEntries(entries []*riotmodel.LeagueEntryDTO, loc string) {
	pipe := p.rdb.Pipeline()
	key := fmt.Sprintf("/entry/%s", loc)
	ctx := context.Background()
	pipe.Expire(ctx, key, p.stgy.LifeTime)
	cmds := make([]*redis.IntCmd, 0, len(entries))
	for _, e := range entries {
		cmds = append(cmds, pipe.HSet(ctx, key, e.SummonerID, e))
	}
	if _, err := pipe.Exec(ctx); err != nil {
		p.logger.Error("redis store entry failed", zap.Error(err))
	}
}

func (p *Pumper) loadEntrie(loc string) {
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
			// assign to localmap
			if _, has := p.entryMap[loc][e.SummonerID]; !has {
				p.entryMap[loc][e.SummonerID] = e
			}
		} // assign to redis
		p.cacheEntries(entries, loc)
	}
	
	// check local && redis map diff
	tmp := make([]*riotmodel.LeagueEntryDTO, 0, p.stgy.MaxSize)
	for k, e := range redisMap {
		if _, has := p.entryMap[loc][k]; !has {
			p.entryMap[loc][k] = e
			tmp = append(tmp, e)
		}
	}
	var max uint = 0
	for k, e := range p.entryMap[loc] {
		if _, has := redisMap[k]; !has {
			tmp = append(tmp, e)
		}
		if max < e.ID {
			max = e.ID
		}
	}
	
	loCode := utils.ConverHostLoCode(loc)
	p.lock.Lock()
	p.entrieIdx[loCode] += (max+1)%(loCode*1e9) + (loCode * 1e9)
	p.lock.Unlock()
	// store diff to db
	p.handleEntries(tmp, loc)
}
