package pumper

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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
	ctx := context.Background()
	key := fmt.Sprintf("/match/%s", loc)
	if _, has := p.matchMap[loc]; !has {
		p.matchMap[loc] = make(map[string]bool)
	}
	// load if redis cache exist
	if size := p.rdb.HLen(ctx, key).Val(); size != 0 {
		keys := p.rdb.HKeys(ctx, key).Val()
		for _, k := range keys {
			p.matchMap[loc][k] = true
		}
	}
	// load if gorm db data exist
	var matches []*riotmodel.MatchDto
	if err := p.db.Where("loc = ?", loc).Find(&matches).Error; err != nil {
		p.logger.Error("load match from gorm db failed", zap.Error(err))
	}
	if len(matches) != 0 {
		for _, m := range matches {
			if _, has := p.matchMap[loc][m.Metadata.MetaMatchID]; !has {
				p.matchMap[loc][m.Metadata.MetaMatchID] = true
			}
		}
	}
}

func (p *Pumper) createMatchID(loCode uint) {
	var (
		url string
	)
	loc, _ := utils.ConvertHostURL(loCode)
	host := utils.ConvertPlatformToHost(loCode)
	p.loadMatch(loc)
	
	// init query val
	startTime := time.Now().AddDate(-1, 0, 0).Unix() // one year ago unix
	endTime := time.Now().Unix()                     // cur time unix
	queryParams := fmt.Sprintf("startTime=%d&endTime=%d&start=0&count=%d",
		startTime, endTime, p.stgy.MaxMatchCount)
	// avoid recursive calls
	copySumnMap := make(map[string]*riotmodel.SummonerDTO)
	for sid, s := range p.sumnMap[loc] {
		copySumnMap[sid] = s
	}
	for sid, summoner := range copySumnMap {
		if summoner.Matches == nil {
			summoner.Matches = make([]*riotmodel.MatchDto, 0, p.stgy.MaxMatchCount)
		}
		// ranker || no rank tier && len(match)<require
		if _, has := p.entryMap[loc][sid]; has || !has && len(summoner.Matches) < p.stgy.MaxMatchCount {
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
}

func (p *Pumper) fetchMatch() {
	var (
		buff      []byte
		err       error
		matches   []*riotmodel.MatchDto
		matchList []string
		url       string
	)
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
				// get match list
				if buff, err = p.fetcher.Get(data.URL); err != nil || buff == nil {
					p.logger.Error(fmt.Sprintf("fetch summoner %s's match failed",
						data.sumn.MetaSummonerID), zap.Error(err))
					if req.Retry < p.stgy.Retry {
						req.Retry++
						p.scheduler.Push(req)
					}
					continue
				}
				if err = json.Unmarshal(buff, &matchList); err != nil {
					p.logger.Error(fmt.Sprintf("unmarshal json to %s failed",
						"sumnDTO"), zap.Error(err))
				}
				oldMatchList := make(map[string]struct{})
				for _, m := range data.sumn.Matches {
					oldMatchList[m.Metadata.MetaMatchID] = struct{}{}
				}
				// init val
				matches = make([]*riotmodel.MatchDto, 0, p.stgy.MaxMatchCount)
				loc := utils.ConverHostLoCode(req.Loc)
				host := utils.ConvertPlatformToHost(loc)
				// fetch match
				for _, matchID := range matchList {
					if _, has := p.matchMap[req.Loc][matchID]; has {
						continue
					}
					if _, has := oldMatchList[matchID]; has {
						continue
					} else {
						url = fmt.Sprintf("%s/lol/match/v5/matches/%s", host, matchID)
						if buff, err = p.fetcher.Get(url); err != nil || buff == nil {
							p.logger.Error(fmt.Sprintf("fetch %s failed", matchID), zap.Error(err))
						}
						var tmp *riotmodel.MatchDto
						if err = json.Unmarshal(buff, &tmp); err != nil {
							p.logger.Error(fmt.Sprintf("unmarshal json to %s failed",
								"MatchDto"), zap.Error(err))
							continue
						}
						matches = append(matches, tmp)
						p.matchMap[req.Loc][matchID] = true
					}
				}
				// deal relation
				summoner := data.sumn
				summoner.Matches = append(summoner.Matches, matches...)
				p.handleMatches(matches, req.Loc)
			}
		}
	}
}
func (p *Pumper) handleMatches(matches []*riotmodel.MatchDto, loc string) {
	loCode := utils.ConverHostLoCode(loc)
	ctx := context.Background()
	key := fmt.Sprintf("/match/%s", loc)
	pipe := p.rdb.Pipeline()
	pipe.Expire(ctx, key, p.stgy.LifeTime)
	cmds := make([]*redis.IntCmd, 0, len(matches))
	tmpSumnList := make([]*riotmodel.SummonerDTO, 0, 10*len(matches))
	
	for _, m := range matches {
		gameId := uint(m.Info.GameID)
		m.ID = loCode*1e11 + gameId
		m.Info.Teams[0].ID = loCode*1e12 + gameId*10 + 1
		m.Info.Teams[1].ID = loCode*1e12 + gameId*10 + 2
		m.Info.Loc = loc
		cmds = append(cmds, pipe.HSet(ctx, key, m.Metadata.MetaMatchID, true))
		for i, par := range m.Info.Participants {
			par.ID = loCode*1e12 + gameId*10 + uint(i)
			if sumn, has := p.sumnMap[loc][par.SummonerId]; has {
				m.Summoners = append(m.Summoners, sumn)
			} else {
				tmp := &riotmodel.SummonerDTO{
					MetaSummonerID: par.SummonerId,
					RevisionDate:   time.Now(),
					Loc:            loc,
				}
				tmpSumnList = append(tmpSumnList, tmp)
				m.Summoners = append(m.Summoners, tmp)
			}
		}
	}
	if len(tmpSumnList) > 0 {
		p.handleSummoner(loc, tmpSumnList...)
	}
	if _, err := pipe.Exec(ctx); err != nil {
		p.logger.Error("redis store match failed", zap.Error(err))
	}
	
	p.out <- &ParseResult{
		Type:  "match",
		Brief: "matches size:" + strconv.Itoa(len(matches)),
		Data:  matches,
	}
}
