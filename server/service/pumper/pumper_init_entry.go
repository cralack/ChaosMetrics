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

func (p *Pumper) InitEntries() {
	exit := make(chan struct{})
	go p.Schedule()
	go p.startTimer()
	for _, loc := range p.stgy.Loc {
		for _, que := range p.stgy.Que {
			// Generate URLs
			go p.createEntriesURL(loc, que)
		}
	}
	// go Watch_ETCD
	go p.fetchDTO()
	// DTO store
	go p.handleResult(exit)
	// blocking until handle final
	<-exit
}

func (p *Pumper) createEntriesURL(loc, que uint) {
	var (
		path string
		tier uint
		rank uint
	)
	locStr, host := utils.ConvertHostURL(loc)
	if p.entriesDic[locStr] == nil {
		p.entriesDic[locStr] = make(map[string]*riotmodel.LeagueEntryDTO)
	}
	queStr := getQueueString(que)
	// generate BEST URL task
	for tier = riotmodel.CHALLENGER; tier <= riotmodel.MASTER; tier++ {
		t, r := ConvertRankToStr(tier, 1)
		
		path = fmt.Sprintf("/lol/league/v4/%sleagues/by-queue/%s",
			strings.ToLower(t), queStr)
		p.scheduler.Push(&scheduler.Task{
			Key:  "best",
			Loc:  locStr,
			Tier: t,
			Rank: r,
			URL:  host + path,
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
				Key:  "mortal",
				Loc:  locStr,
				Tier: t,
				Rank: r,
				URL:  host + path,
			})
			
		}
	}
}

func (p *Pumper) fetchDTO() {
	
	var (
		page    int
		buff    []byte
		err     error
		list    riotmodel.LeagueListDTO
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
		p.logger.Info(fmt.Sprintf("fetching request URL:%s", req.URL))
		// fetch and parse
		switch req.Key {
		case "best":
			// rate limiter
			<-p.rater
			// api: /lol/league/v4/{BEST}leagues/by-queue/{queue}
			if buff, err = p.fetcher.Get(req.URL); err != nil {
				p.logger.Error(fmt.Sprintf("fetch %s %s failed", req.Tier, req.Rank),
					zap.Error(err))
			}
			if err = json.Unmarshal(buff, &list); err != nil {
				p.logger.Error(fmt.Sprintf("unmarshal json to %s failed",
					"LeagueListDTO"), zap.Error(err))
			}
			entries = list.Entries
			p.handleEntry(entries, req, 1)
			p.cacheEntries(entries, req)
			// check oversize && chunk
		
		case "mortal":
			page = 0
			for {
				page++
				// rate limiter
				<-p.rater
				// api: /lol/league/v4/entries/{queue}/{tier}/{division}
				if buff, err = p.fetcher.Get(fmt.Sprintf("%s?page=%s",
					req.URL, strconv.Itoa(page))); err != nil {
					p.logger.Error(fmt.Sprintf("fetch %s %s failed", req.Tier, req.Rank),
						zap.Error(err))
				}
				if err = json.Unmarshal(buff, &entries); err != nil {
					p.logger.Error(fmt.Sprintf("unmarshal json to %s failed",
						"LeagueItemDTO"), zap.Error(err))
				}
				// send finish signal
				if len(entries) == 0 {
					p.logger.Info(fmt.Sprintf("all %s %s data fetch done at page %d",
						req.Tier, req.Rank, page))
					if req.Tier == endTier && req.Rank == endRank {
						p.out <- &ParseResult{
							Type:  "entry",
							Brief: "finish",
							Page:  0,
							Data:  nil,
						}
					}
					break
				}
				p.handleEntry(entries, req, page)
				p.cacheEntries(entries, req)
			}
		}
	}
}

func (p *Pumper) handleEntry(entries []*riotmodel.LeagueEntryDTO, req *scheduler.Task, page int, ) {
	tmp := make([]*riotmodel.LeagueEntryDTO, 0, 500)
	loc := utils.ConverHostLoCode(req.Loc)
	
	p.lock.Lock()
	defer p.lock.Unlock()
	for _, entry := range entries {
		if _, has := p.entriesDic[req.Loc][entry.SummonerID]; !has {
			p.entriesDic[req.Loc][entry.SummonerID] = entry
			entry.Loc = req.Loc
			entry.Tier = req.Tier
			entry.ID = p.entrieIdx[loc] + loc*1e9 // idx+offset
			p.entrieIdx[loc]++
			tmp = append(tmp, entry)
		}
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
	for i, chunk := range chunks {
		p.out <- &ParseResult{
			Type:  "entry",
			Brief: req.Tier + " " + req.Rank,
			Page:  i + page,
			Data:  chunk,
		}
	}
}

func (p *Pumper) cacheEntries(entries []*riotmodel.LeagueEntryDTO, req *scheduler.Task) {
	pipe := p.rdb.Pipeline()
	key := fmt.Sprintf("/entry/%s", req.Loc)
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

func (p *Pumper) handleResult(exit chan struct{}) {
	for result := range p.out {
		switch result.Type {
		case "entry":
			if result.Data == nil {
				p.logger.Info("all result handle done")
				exit <- struct{}{}
				break
			}
			data := result.Data.([]*riotmodel.LeagueEntryDTO)
			size := len(data)
			if err := p.db.Create(data).Error; err != nil {
				p.logger.Error("riot model store failed", zap.Error(err))
			} else {
				p.logger.Info(fmt.Sprintf("%s's entries store succeed %d at page %02d",
					data[0].Tier+" "+data[0].Rank, size, result.Page))
			}
			
		}
		
	}
}
