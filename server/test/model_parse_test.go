package test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"testing"
	
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func Test_parse_summoners(t *testing.T) {
	// fetching remote JSON data (3~5 seconds per request)
	// url := "https://tw2.api.riotgames.com/lol/summoner/v4/summoners/by-name/Mudife"
	// buff, err := f.Get(url)
	buff, err := os.ReadFile(path + "summoners.txt")
	if err != nil {
		t.Fatal(err)
	}
	
	// 解析 JSON 数据
	var summoners []*riotmodel.SummonerDTO
	err = json.Unmarshal(buff, &summoners)
	if err != nil {
		fmt.Println("解析失败：", err)
		return
	}
	
	// 打印解析结果
	for _, summoner := range summoners {
		fmt.Println("ID:", summoner.ID)
		fmt.Println("AccountID:", summoner.AccountID)
		fmt.Println("PUUID:", summoner.PUUID)
		fmt.Println("Name:", summoner.Name)
		fmt.Println("ProfileIconID:", summoner.ProfileIconID)
		fmt.Println("RevisionDate:", summoner.RevisionDate)
		fmt.Println("SummonerLevel:", summoner.SummonerLevel)
		fmt.Println()
	}
	
	// gorm create
	if err := db.Create(summoners).Error; err != nil {
		logger.Info("db create summoners failed")
	}
	// gorm read
	var tar []*riotmodel.SummonerDTO
	if err = db.Find(&tar).Error; err != nil {
		logger.Error("db ")
	}
	
	// redis create
	pipe := rdb.Pipeline()
	ctx := context.Background()
	key := "/summoner/tw2"
	cmds := make([]*redis.IntCmd, 0, len(summoners))
	for _, s := range summoners {
		cmds = append(cmds, pipe.HSet(ctx, key, s.MetaSummonerID))
	}
	if err = rdb.HSet(ctx, key, summoners[0].MetaSummonerID, summoners[0]).Err(); err != nil {
		logger.Info("")
	}
	if _, err = pipe.Exec(ctx); err != nil {
		logger.Error("redis create summoners failed")
	}
	// redis read
	redisMap := make(map[string]*riotmodel.SummonerDTO)
	kvmap := rdb.HGetAll(ctx, key).Val()
	for k, v := range kvmap {
		var tmp riotmodel.SummonerDTO
		if err := json.Unmarshal([]byte(v), &tmp); err != nil {
			logger.Error("load entry form redis cache failed", zap.Error(err))
		} else {
			redisMap[k] = &tmp
		}
	}
}

func Test_parse_championr_rotation(t *testing.T) {
	// fetching remote JSON data (3~5 seconds per request)
	url := "https://tw2.api.riotgames.com/lol/platform/v3/champion-rotations"
	buff, err := f.Get(url)
	
	// load local json data
	// buff, err := os.ReadFile(path + "championr_rotation.txt")
	if err != nil {
		t.Fatal(err)
	}
	var res riotmodel.ChampionRotationDTO
	err = json.Unmarshal(buff, &res)
	if err != nil {
		t.Fatal(err)
	}

	freeForNew := res.FreeChampionIdsForNewPlayers
	t.Log(freeForNew)
}

func Test_parse_championr_mastery(t *testing.T) {
	// fetching remote JSON data (3~5 seconds per request)
	// url := "https://tw2.api.riotgames.com/lol/champion-mastery/v4/champion-masteries/by-puuid/F4fFtqehQLBj8U5sKBZF--k-7akbtb1IX790lRd4whPI4pXDAuVyfswHetg2lz_kMe2NJ0gUo5EIig/top"
	// buff, err := f.Get(url)
	buff, err := os.ReadFile(path + "championr_mastery.txt")
	if err != nil {
		t.Fatal(err)
	}
	var res []riotmodel.ChampionMasteryDto
	err = json.Unmarshal(buff, &res)
	if err != nil {
		t.Fatal(err)
	}

	for _, mastery := range res {
		fmt.Println("Champion ID:", mastery.ChampionID)
		fmt.Println("Champion Level:", mastery.ChampionLevel)
		fmt.Println("Champion Points:", mastery.ChampionPoints)
		fmt.Println("Last Play Time:", mastery.LastPlayTime)
		fmt.Println("Champion Points Since Last Level:", mastery.ChampionPointsSinceLastLevel)
		fmt.Println("Champion Points Until Next Level:", mastery.ChampionPointsUntilNextLevel)
		fmt.Println("Chest Granted:", mastery.ChestGranted)
		fmt.Println("Tokens Earned:", mastery.TokensEarned)
		fmt.Println("Summoner ID:", mastery.SummonerID)
		fmt.Println()
	}
}

func Test_parse_match(t *testing.T) {
	// fetching remote JSON data (3~5 seconds per request)
	// url := "https://sea.api.riotgames.com/lol/match/v5/matches/TW2_81882122"
	// buff, err := f.Get(url)

	// load local json data
	buff, err := os.ReadFile(path + "match.txt")
	if err != nil {
		t.Fatal(err)
	}
	// parse
	var res *riotmodel.MatchDto
	err = json.Unmarshal(buff, &res)
	if err != nil {
		t.Fatal(err)
	}
	// check data
	fmt.Println("Metadata Data Version:",
		res.Metadata.DataVersion)
	fmt.Println("Game Creation Time:",
		res.Info.GameCreation)
	fmt.Println("Game Version:",
		res.Info.GameVersion)
	fmt.Println("Win Status of Team 1:",
		res.Info.Teams[0].Win)
	fmt.Println("Baron Kills for Team 1:",
		res.Info.Teams[0].Objectives.Baron.Kills)

	participants := res.Info.Participants

	// Find the participant with FirstBloodKill
	var idx int
	for i, p := range participants {
		if p.FirstBloodKill {
			fmt.Println("Champion Name of First Blood Kill:",
				p.ChampionName)
			idx = i
			break
		}
	}
	
	// Retrieve information for the participant with FirstBloodKill
	player := participants[idx]
	fmt.Println("Kills/Deaths/Assists of First Blood Player:",
		player.Kills, player.Deaths, player.Assists)
	fmt.Println("Total Damage Dealt to Champions by First Blood Player:",
		player.TotalDamageDealtToChampions)
	fmt.Println("First Blood Player's match ID:",
		player.MetaMatchID)
	
	// gorm create
	if err := db.Create(res).Error; err != nil {
		logger.Info("db create summoners failed")
	}
	// gorm read
	var tar *riotmodel.MatchDto
	if err = db.Where("meta_match_id = ?", "TW2_81882122").Find(&tar).Error; err != nil {
		logger.Error("db ")
	} else {
		logger.Info(fmt.Sprintf("res == tar:%v", res.Metadata.MetaMatchID == tar.Metadata.MetaMatchID))
	}
	
	// redis create
	ctx := context.Background()
	key := "/match/tw2"
	if err := rdb.HSet(ctx, key, res.Metadata.MetaMatchID, res).Err(); err != nil {
		logger.Error("redis create match failed")
	}
	
	result := rdb.HGet(ctx, key, tar.Metadata.MetaMatchID).Val()
	var tar2 *riotmodel.MatchDto
	if err := json.Unmarshal([]byte(result), &tar2); err != nil {
		logger.Error("redis read match failed")
	}
	// redis read
	redisMap := make(map[string]*riotmodel.SummonerDTO)
	kvmap := rdb.HGetAll(ctx, key).Val()
	for k, v := range kvmap {
		var tmp riotmodel.SummonerDTO
		if err := json.Unmarshal([]byte(v), &tmp); err != nil {
			logger.Error("load entry form redis cache failed", zap.Error(err))
		} else {
			redisMap[k] = &tmp
		}
		
	}
}

func Test_parse_league(t *testing.T) {
	// CHALLENGER
	// url := "https://tw2.api.riotgames.com/lol/league/v4/challengerleagues/by-queue/RANKED_SOLO_5x5"
	// GRANDMASTER
	// url := "https://tw2.api.riotgames.com/lol/league/v4/grandmasterleagues/by-queue/RANKED_SOLO_5x5"
	// MASTER
	// url := "https://tw2.api.riotgames.com/lol/league/v4/masterleagues/by-queue/RANKED_SOLO_5x5"
	url := "https://tw2.api.riotgames.com/lol/league/v4/masterleagues/by-queue/RANKED_FLEX_SR"
	buff, err := f.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	// parse
	var res riotmodel.LeagueListDTO
	err = json.Unmarshal(buff, &res)
	if err != nil {
		t.Fatal(err)
	}
	entries := res.Entries
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].LeaguePoints > entries[j].LeaguePoints
	})
	n := 10
	for i := 0; i < n; i++ {
		fmt.Printf("summoner %s's LP is %d\n", entries[i].SummonerName, entries[i].LeaguePoints)
	}
}

func Test_parse_mortal(t *testing.T) {
	url := "https://tw2.api.riotgames.com/lol/league/v4/entries/RANKED_SOLO_5x5/DIAMOND/I?page=1"
	buff, err := f.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	// parse
	var res []*riotmodel.LeagueEntryDTO
	err = json.Unmarshal(buff, &res)
	if err != nil {
		t.Fatal(err)
	}
	
}