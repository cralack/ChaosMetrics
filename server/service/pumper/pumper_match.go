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
)

type matchTask struct {
	sumn *riotmodel.SummonerDTO
	URL  string
}

func (p *Pumper) UpdateMatch(exit chan struct{}) {
	for _, loc := range p.stgy.Loc {
		go p.createMatchID(loc)
	}
	go p.fetchMatch()
	<-exit
}

func (p *Pumper) loadMatch(loc string) {
	var (
		matches  []*riotmodel.MatchDto
		redisMap map[string]bool
		size     int64
		err      error
	)
	if _, has := p.matchMap[loc]; !has {
		p.matchMap[loc] = make(map[string]bool)
	}
	// load if gorm db data exist
	if err = p.db.Where("loc = ?", loc).Find(&matches).Error; err != nil {
		p.logger.Error("load match from gorm db failed", zap.Error(err))
	}
	if len(matches) != 0 {
		for _, m := range matches {
			// assign to local map
			if _, has := p.matchMap[loc][m.Metadata.MetaMatchID]; !has {
				p.matchMap[loc][m.Metadata.MetaMatchID] = true
			}
		}
	}
	ctx := context.Background()
	key := fmt.Sprintf("/match/%s", loc)
	
	// load if redis cache exist
	if size = p.rdb.HLen(ctx, key).Val(); size != 0 {
		redisMap = make(map[string]bool, size)
		keys := p.rdb.HKeys(ctx, key).Val()
		dels := make([]string, 0, size)
		for _, k := range keys {
			if _, has := p.matchMap[loc][k]; !has {
				dels = append(dels, k)
			} else {
				redisMap[k] = true
			}
		}
		// sync db to cache
		p.rdb.HDel(ctx, key, dels...)
	}
	pipe := p.rdb.Pipeline()
	cmds := make([]*redis.IntCmd, 0, len(matches))
	for _, m := range matches {
		if _, has := redisMap[m.Metadata.MetaMatchID]; !has {
			cmds = append(cmds, pipe.HSet(ctx, key, m.Metadata.MetaMatchID, true))
		}
	}
	if _, err = pipe.Exec(ctx); err != nil {
		p.logger.Error("sync match form db to cache failed", zap.Error(err))
	}
}

func (p *Pumper) createMatchID(loCode uint) {
	var (
		url      string
		sid      string
		summoner *riotmodel.SummonerDTO
		has      bool
	)
	loc, _ := utils.ConvertHostURL(loCode)
	host := utils.ConvertPlatformToHost(loCode)
	p.loadMatch(loc)
	
	// init query val
	startTime := time.Now().AddDate(-1, 0, 0).Unix() // one year ago unix
	endTime := time.Now().Unix()                     // cur time unix
	queryParams := fmt.Sprintf("startTime=%d&endTime=%d&start=0&count=%d",
		startTime, endTime, p.stgy.MaxMatchCount)
	
	for sid, summoner = range p.sumnMap[loc] {
		matchList := utils.ConvertStrToSlice(summoner.Matches)
		// ranker || no rank tier && len(match)<require
		if _, has = p.entryMap[loc][sid]; has || !has && len(matchList) < p.stgy.MaxMatchCount {
			url = fmt.Sprintf("%s/lol/match/v5/matches/by-puuid/%s/ids?%s",
				host, summoner.PUUID, queryParams)
			p.scheduler.Push(&scheduler.Task{
				Key:   "match",
				Loc:   loc,
				Retry: 0,
				Data: &matchTask{
					sumn: summoner,
					URL:  url,
				},
			})
		}
	}
	
	// finish signal
	p.scheduler.Push(&scheduler.Task{
		Key:  "match",
		Loc:  loc,
		Data: nil,
	})
	return
}

func (p *Pumper) fetchMatch() {
	var (
		buff         []byte
		cnt          int
		err          error
		matches      []*riotmodel.MatchDto
		curMatchList []string
		url          string
	)
	defer func() {
		if err := recover(); err != nil {
			p.logger.Panic("fetcher panic",
				zap.Any("err", err),
				zap.String("stack", string(debug.Stack())))
		}
	}()
	
	for {
		req := p.scheduler.Pull()
		switch req.Key {
		case "match":
			if req.Data == nil {
				// send finish signal
				p.out <- &ParseResult{
					Type:  "finish",
					Brief: "match",
					Data:  nil,
				}
				// *need release scheduler resource*
				return
			} else {
				data := req.Data.(*matchTask)
				// get old & cur match list
				if buff, err = p.fetcher.Get(data.URL); err != nil || buff == nil {
					p.logger.Error(fmt.Sprintf("fetch summoner %s's match failed",
						data.sumn.MetaSummonerID), zap.Error(err))
					if req.Retry < p.stgy.Retry {
						req.Retry++
						p.scheduler.Push(req)
					}
					continue
				}
				if err = json.Unmarshal(buff, &curMatchList); err != nil {
					p.logger.Error(fmt.Sprintf("unmarshal json to %s failed",
						"sumnDTO"), zap.Error(err))
				}
				// get old match list
				oldMatchList := make(map[string]struct{})
				for _, matchID := range utils.ConvertStrToSlice(data.sumn.Matches) {
					oldMatchList[matchID] = struct{}{}
				}
				// update summoner's match list
				summoner := data.sumn
				summoner.Matches = utils.ConvertSliceToStr(curMatchList)
				cnt++
				p.handleSummoner(req.Loc, summoner)
				// init val
				matches = make([]*riotmodel.MatchDto, 0, p.stgy.MaxMatchCount)
				loc := utils.ConverHostLoCode(req.Loc)
				host := utils.ConvertPlatformToHost(loc)
				// fetch match
				for _, matchID := range curMatchList {
					if _, has := p.matchMap[req.Loc][matchID]; has {
						continue
					}
					if _, has := oldMatchList[matchID]; has {
						continue
					} else {
						url = fmt.Sprintf("%s/lol/match/v5/matches/%s", host, matchID)
						if buff, err = p.fetcher.Get(url); err != nil || buff == nil {
							p.logger.Error(fmt.Sprintf("fetch %s failed", matchID), zap.Error(err))
							if req.Retry < p.stgy.Retry {
								req.Retry++
								p.scheduler.Push(req)
							}
							continue
						}
						var tmp *riotmodel.MatchDto
						if err = json.Unmarshal(buff, &tmp); err != nil {
							p.logger.Error(fmt.Sprintf("unmarshal json to %s failed",
								"MatchDto"), zap.Error(err))
							continue
						}
						// remake?
						if tmp.Info.GameDuration == 0 {
							p.logger.Info(summoner.Name + "'s remake match")
						} else {
							matches = append(matches, tmp)
						}
						p.matchMap[req.Loc][matchID] = true
					}
				}
				p.logger.Info(fmt.Sprintf(
					"updating %s's match list @ %d,store %d matches", summoner.Name, cnt, len(matches)))
				if len(matches) == 0 {
					continue
				}
				p.handleMatches(matches, req.Loc, summoner.Name)
			}
		}
	}
}

func (p *Pumper) handleMatches(matches []*riotmodel.MatchDto, loc, sName string) {
	if len(matches) == 0 {
		return
	}
	loCode := utils.ConverHostLoCode(loc)
	ctx := context.Background()
	key := fmt.Sprintf("/match/%s", loc)
	pipe := p.rdb.Pipeline()
	pipe.Expire(ctx, key, p.stgy.LifeTime)
	cmds := make([]*redis.IntCmd, 0, len(matches))
	
	for _, m := range matches {
		if m == nil || m.Info == nil || len(m.Info.Teams) != 2 ||
			len(m.Info.Participants) < 8 {
			p.logger.Error(sName + "'s wrong match")
			continue
		}
		gameId := uint(m.Info.GameID)
		m.ID = loCode*1e11 + gameId
		m.Info.Teams[0].ID = loCode*1e12 + gameId*10 + 1
		m.Info.Teams[1].ID = loCode*1e12 + gameId*10 + 2
		m.Info.Loc = loc
		cmds = append(cmds, pipe.HSet(ctx, key, m.Metadata.MetaMatchID, true))
		for i, par := range m.Info.Participants {
			par.ID = loCode*1e12 + gameId*10 + uint(i)
		}
	}
	if _, err := pipe.Exec(ctx); err != nil {
		p.logger.Error("redis store match failed", zap.Error(err))
	}
	
	p.out <- &ParseResult{
		Type:  "match",
		Brief: sName,
		Data:  matches,
	}
	return
}

