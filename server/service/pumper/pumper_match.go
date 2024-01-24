package pumper

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/cralack/ChaosMetrics/server/utils/scheduler"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type matchTask struct {
	sumn *riotmodel.SummonerDTO
}

func (p *Pumper) UpdateMatch(exit chan struct{}) {
	for _, loc := range p.stgy.Loc {
		go p.createMatchListURL(loc)
	}
	<-exit
}

func (p *Pumper) loadMatch(loc string) {
	var (
		matches  []*riotmodel.MatchDB
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
			if _, has := p.matchMap[loc][m.MetaMatchID]; !has {
				p.matchMap[loc][m.MetaMatchID] = true
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
		if _, has := redisMap[m.MetaMatchID]; !has {
			cmds = append(cmds, pipe.HSet(ctx, key, m.MetaMatchID, true))
		}
	}
	if _, err = pipe.Exec(ctx); err != nil {
		p.logger.Error("sync match form db to cache failed", zap.Error(err))
	}
	matches = nil
	return
}

func (p *Pumper) createMatchListURL(loCode riotmodel.LOCATION) {
	var (
		url      string
		sid      string
		summoner *riotmodel.SummonerDTO
		has      bool
		count    int
	)
	loc, _ := utils.ConvertHostURL(loCode)
	region := utils.ConvertLocToRegion(loCode)
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
				region, summoner.PUUID, queryParams)
			p.scheduler.Push(&scheduler.Task{
				Type: matchTypeKey,
				Loc:  loc,
				URL:  url,
				Data: &matchTask{
					sumn: summoner,
				},
			})
			count++
		}
	}

	// finish signal
	p.scheduler.Push(&scheduler.Task{
		Type: matchTypeKey,
		Loc:  loc,
		Data: nil,
	})
	return
}

func (p *Pumper) FetchMatchByID(req *scheduler.Task, host, matchID string) (res *riotmodel.MatchDB) {
	var (
		buff    []byte
		url     string
		err     error
		match   *riotmodel.MatchDTO
		matchTL *riotmodel.MatchTimelineDTO
	)

	sumName := req.Data.(*matchTask).sumn.Name
	// fetch match
	url = fmt.Sprintf("%s/lol/match/v5/matches/%s", host, matchID)
	if buff, err = p.fetcher.Get(url); err != nil || len(buff) < 1000 {
		p.logger.Error(fmt.Sprintf("fetch %s's match %s failed",
			sumName, matchID), zap.Error(err))
		if req.Retry < p.stgy.Retry {
			req.Retry++
			p.scheduler.Push(req)
		}
		return
	}
	if err = json.Unmarshal(buff, &match); err != nil {
		p.logger.Error(fmt.Sprintf("unmarshal %s's match %s json failed",
			sumName, matchID), zap.Error(err))
		return
	}
	// remake || bot game
	if match.Info.GameDuration == 0 || len(match.Metadata.Participants) <= 5 ||
		// won't parse cherry mode for now
		match.Info.GameMode == "CHERRY" {
		return
	}
	// fetch match timeline
	url = fmt.Sprintf("%s/lol/match/v5/matches/%s/timeline", host, matchID)
	if buff, err = p.fetcher.Get(url); err != nil || len(buff) < 1000 {
		p.logger.Error(fmt.Sprintf("fetch %s's match timeline %s failed",
			sumName, matchID), zap.Error(err))
		if req.Retry < p.stgy.Retry {
			req.Retry++
			p.scheduler.Push(req)
		}
		return
	}
	if err = json.Unmarshal(buff, &matchTL); err != nil {
		p.logger.Error(fmt.Sprintf("unmarshal %s's match %s json failed",
			sumName, matchID), zap.Error(err))
		return
	}

	res = &riotmodel.MatchDB{
		Analyzed:       false,
		MetaMatchID:    match.Metadata.MatchID,
		Loc:            strings.ToLower(match.Info.PlatformID),
		GameCreation:   match.Info.GameCreation,
		GameDuration:   match.Info.GameDuration,
		GameMode:       match.Info.GameMode,
		GameVersion:    match.Info.GameVersion,
		MapID:          match.Info.MapID,
		QueueID:        match.Info.QueueID,
		TournamentCode: match.Info.TournamentCode,
	}
	if err = res.ParseClassicAndARAM(match, matchTL); err != nil {
		p.logger.Error("parse match failed", zap.Error(err))
		return nil
	}
	return
}

func (p *Pumper) handleMatches(matches []*riotmodel.MatchDB, sName string) {
	if len(matches) == 0 {
		return
	}

	ctx := context.Background()
	pipe := p.rdb.Pipeline()
	// pipe.Expire(ctx, key, p.stgy.LifeTime)
	cmds := make([]*redis.IntCmd, 0, len(matches))

	for _, m := range matches {
		loCode := utils.ConverHostLoCode(m.Loc)
		key := fmt.Sprintf("/match/%s", m.Loc)
		cmds = append(cmds, pipe.HSet(ctx, key, m.MetaMatchID, true))
		var gameId uint
		if id, err := strconv.Atoi(m.MetaMatchID[4:]); err != nil {
			p.logger.Error("transfer num failed", zap.Error(err))
			continue
		} else {
			gameId = uint(id)
		}
		m.ID = gameId*1e2 + uint(loCode)
		// m.Loc = loc
		for i, par := range m.Participants {
			par.ID = gameId*1e3 + uint(i)*1e2 + uint(loCode)
		}
	}
	if _, err := pipe.Exec(ctx); err != nil {
		p.logger.Error("redis store match failed", zap.Error(err))
	}

	var chunks [][]*riotmodel.MatchDB
	totalSize := len(matches)
	chunkSize := 5
	if totalSize < chunkSize {
		p.out <- &DBResult{
			Type:  matchTypeKey,
			Brief: sName,
			Data:  matches,
		}
		return
	} else {
		for i := 0; i < totalSize; i += chunkSize {
			end := i + chunkSize
			if end > totalSize {
				end = totalSize
			}
			chunks = append(chunks, matches[i:end])
		}
	}

	for _, chunk := range chunks {
		p.out <- &DBResult{
			Type:  matchTypeKey,
			Brief: sName,
			Data:  chunk,
		}
	}
	return
}

func (p *Pumper) FetchMatchByName(summonerName string, loc riotmodel.LOCATION) error {
	var (
		puuid  string
		region string
		url    string
		locStr string
		sumn   *riotmodel.SummonerDTO
	)
	locStr, _ = utils.ConvertHostURL(loc)
	region = utils.ConvertLocToRegion(loc)
	sumn = p.LoadSingleSummoner(summonerName, locStr)
	puuid = sumn.PUUID
	url = fmt.Sprintf("")
	startTime := time.Now().AddDate(-1, 0, 0).Unix() // one year ago unix
	endTime := time.Now().Unix()                     // cur time unix
	queryParams := fmt.Sprintf("startTime=%d&endTime=%d&start=0&count=%d",
		startTime, endTime, p.stgy.MaxMatchCount)

	url = fmt.Sprintf("%s/lol/match/v5/matches/by-puuid/%s/ids?%s",
		region, puuid, queryParams)

	p.scheduler.Push(&scheduler.Task{
		Type: matchTypeKey,
		Loc:  locStr,
		URL:  url,
		Data: &matchTask{
			sumn: sumn,
		},
		Priority: true,
	})
	return nil
}
