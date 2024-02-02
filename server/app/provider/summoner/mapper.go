package summoner

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

func (s *SumonerService) HandleSummoner(src *riotmodel.SummonerDTO) (des *response.SummonerDTO) {
	var (
		key       string
		field     string
		buff      string
		matchList []string
		err       error
		ctx       = context.Background()
	)
	if src == nil {
		return nil
	}
	des = &response.SummonerDTO{
		Name:          src.Name,
		Loc:           src.Loc,
		ProfileIconID: src.ProfileIconID,
		SummonerLevel: src.SummonerLevel,
	}
	// load entry from redis
	{
		key = fmt.Sprintf("/entry/%s", src.Loc)
		field = fmt.Sprintf("%s@RANKED_SOLO_5x5", src.Name)
		buff = s.rdb.HGet(ctx, key, field).Val()
		entry := &riotmodel.LeagueEntryDTO{}
		if err = json.Unmarshal([]byte(buff), &entry); err == nil && entry.SummonerName == src.Name {
			des.SoloEntry = ConvertEntry(entry)
		}
		field = fmt.Sprintf("%s@RANKED_FLEX_SR", src.Name)
		buff = s.rdb.HGet(ctx, key, field).Val()
		entry = &riotmodel.LeagueEntryDTO{}
		if err = json.Unmarshal([]byte(buff), &entry); err == nil && entry.SummonerName == src.Name {
			des.FlexEntry = ConvertEntry(entry)
		}
	}

	matchList = utils.ConvertStrToSlice(src.Matches)
	des.Matches = make([]*response.MatchDTO, 0, len(matchList))
	// load match from redis
	key = fmt.Sprintf("/matchDTO/%s", src.Loc)
	{
		for _, m := range matchList {
			matchDTO := &response.MatchDTO{}
			buff = s.rdb.HGet(ctx, key, m).Val()
			if err = json.Unmarshal([]byte(buff), &matchDTO); err == nil && matchDTO.MatchID == m {
				des.Matches = append(des.Matches, matchDTO)
			}
		}
		if len(des.Matches) == len(matchList) {
			return
		}
	}
	// load match from db
	{
		loCode := utils.ConvertLocStrToLocation(src.Loc)
		pipe := s.rdb.Pipeline()
		pipe.Expire(ctx, key, time.Hour*24*7)
		cmds := make([]*redis.IntCmd, 0, len(matchList))
		for _, m := range matchList {
			id, _ := strconv.Atoi(m[4:])
			id = id*1e2 + int(loCode)
			var (
				match    *riotmodel.MatchDB
				matchDTO *response.MatchDTO
			)
			err = s.db.Where("id=?", id).Preload(clause.Associations).Find(&match).Error
			if err == nil && match.MetaMatchID == m {
				if matchDTO = ConvertMatchesToDTO(match); matchDTO != nil {
					cmds = append(cmds, pipe.HSet(ctx, key, m, matchDTO))
					des.Matches = append(des.Matches, matchDTO)
				}
			}
		}
		if len(des.Matches) != len(matchList) {
			return nil
		}
		if _, err = pipe.Exec(ctx); err != nil {
			s.logger.Error("redis store matchDTO failed", zap.Error(err))
		}
	}
	return
}

func ConvertMatchesToDTO(src *riotmodel.MatchDB) (des *response.MatchDTO) {
	if src == nil {
		return nil
	}
	des = &response.MatchDTO{
		MatchID:      src.MetaMatchID,
		QueueID:      src.QueueID,
		GameCreation: src.GameCreation.Unix(),
		GameDuration: src.GameDuration,
		Participants: make([]*response.Participant, len(src.Participants)),
	}
	for i, p := range src.Participants {
		des.Participants[i] = &response.Participant{
			SummonerName: p.SummonerName,
			ChampionName: p.ChampionName,
			Kills:        p.Kills,
			Deaths:       p.Deaths,
			Assists:      p.Assists,
			Item0:        p.Item0,
			Item1:        p.Item1,
			Item2:        p.Item2,
			Item3:        p.Item3,
			Item4:        p.Item4,
			Item5:        p.Item5,
			Item6:        p.Item6,
			KDA:          p.KDA,
			KP:           p.KP,
			PentaKills:   p.PentaKills,
			QuadraKills:  p.QuadraKills,
			TripleKills:  p.TripleKills,
			Summoner1Id:  p.Summoner1Id,
			Summoner2Id:  p.Summoner2Id,
			TeamId:       p.TeamId,
			DamageDealt:  p.DamageDealt,
			DamageToken:  p.DamageToken,
			SkillBuild:   p.Build.Skill,
			ItemBuild:    p.Build.Item,
			PerkBuild:    p.Build.Perk,
		}
	}
	return
}

func ConvertEntry(src *riotmodel.LeagueEntryDTO) (des *response.EntryDTO) {
	return &response.EntryDTO{
		QueType:      src.QueType,
		Tier:         src.Tier,
		Rank:         src.Rank,
		LeaguePoints: src.LeaguePoints,
		Wins:         src.Wins,
		Losses:       src.Losses,
	}
}
