package test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"testing"
	"time"
	"unsafe"

	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func Test_parse_summoners(t *testing.T) {
	// fetching remote JSON data (3~5 seconds per request)
	// url := "https://tw2.api.riotgames.com/lol/summoner/v4/summoners/by-name/Mudife"
	// buff, err := f.Get(url)
	buff, err := os.ReadFile(path + "summoners.json")
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
		fmt.Println("ProfileIconID:", summoner.ProfileIconID)
		fmt.Println("RevisionDate:", summoner.RevisionDate)
		fmt.Println("SummonerLevel:", summoner.SummonerLevel)
		fmt.Println()
	}

	// gorm create
	if err := db.Create(summoners).Error; err != nil {
		logger.Debug("orm create summoners failed")
	}
	// gorm read
	var tar []*riotmodel.SummonerDTO
	if err = db.Find(&tar).Error; err != nil {
		logger.Error("orm ")
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
		logger.Debug("")
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

func Test_parse_champion_rotation(t *testing.T) {
	// fetching remote JSON data (3~5 seconds per request)
	url := "https://tw2.api.riotgames.com/lol/platform/v3/champion-rotations"
	buff, err := f.Get(url)

	// load local json data
	// buff, err := os.ReadFile(path + "championr_rotation.json")
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

func Test_parse_champion_mastery(t *testing.T) {
	// fetching remote JSON data (3~5 seconds per request)
	// url := "https://tw2.api.riotgames.com/lol/champion-mastery/v4/champion-masteries/by-puuid/F4fFtqehQLBj8U5sKBZF--k-7akbtb1IX790lRd4whPI4pXDAuVyfswHetg2lz_kMe2NJ0gUo5EIig/top"
	// buff, err := f.Get(url)
	buff, err := os.ReadFile(path + "championr_mastery.json")
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
	buff, err := os.ReadFile(path + "match.json")
	if err != nil {
		t.Fatal(err)
	}
	// parse
	var res *riotmodel.MatchDTO
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
		// fmt.Printf("%s 's TimeCCingOthers:%d\n\r", p.ChampionName, p.TimeCCingOthers)
		// fmt.Printf("%s 's TotalTimeCCDealt:%d\n\r", p.ChampionName, p.TotalTimeCCDealt)
		if p.FirstBloodKill {
			fmt.Println("Champion RiotName of First Blood Kill:",
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
}

func Test_parse_match_list(t *testing.T) {
	const (
		maxMatch = 20
		loc      = riotmodel.TW2
	)
	var (
		buff []byte
		err  error
		list []string
	)
	puuid := "6RtCOQdb0rlWO0S714_nYds_xDw2-bwrB8IsVzJbAQi8uBfosVT5UyfayA9oirdE5pCFMEFB6TkAlA"
	region := utils.ConvertLocationToRegionHost(loc)
	startTime := time.Now().AddDate(-1, 0, 0).Unix() // one year ago unix
	endTime := time.Now().Unix()                     // cur time unix
	queryParams := fmt.Sprintf("startTime=%d&endTime=%d&start=0&count=%d", startTime, endTime, maxMatch)
	url := fmt.Sprintf("%s/lol/match/v5/matches/by-puuid/%s/ids?%s", region, puuid, queryParams)
	if buff, err = f.Get(url); err != nil {
		t.Log(err)
	}

	if err = json.Unmarshal(buff, &list); err != nil {
		t.Log(err)
	}
	t.Log(list)
	ques := make([]int, 0, 20)
	for _, matchId := range list {
		url := fmt.Sprintf("%s/lol/match/v5/matches/%s", region, matchId)
		if buff, err = f.Get(url); err != nil {
			t.Log(err)
		}
		var match *riotmodel.MatchDTO
		if err = json.Unmarshal(buff, &match); err != nil {
			t.Log(err)
		}
		ques = append(ques, match.Info.QueueID)
	}
	t.Log(ques)
}

func Test_parse_matchs(t *testing.T) {
	var (
		buff  []byte
		err   error
		match *riotmodel.MatchDTO
	)
	matchList := []string{"TW2_97186217", "TW2_97161015", "TW2_97131970", "TW2_97119501", "TW2_97106562", "TW2_78013194", "TW2_71558421", "TW2_59549103", "TW2_59539001", "TW2_59527321", "TW2_59516488", "TW2_59509451", "TW2_59502083", "TW2_42926206", "TW2_30891005", "TW2_30846799", "TW2_25739350", "TW2_25511607", "TW2_25501389", "TW2_25454450"}
	for _, matchId := range matchList {
		url := "https://sea.api.riotgames.com/lol/match/v5/matches/" + matchId
		buff, err = f.Get(url)
		if err != nil {
			t.Log(err)
		}
		if err = json.Unmarshal(buff, &match); err != nil {
			t.Log(err)
		}
	}
}

func Test_parse_matchLine(t *testing.T) {
	// fetching remote JSON data (3~5 seconds per request)
	// url := "https://sea.api.riotgames.com/lol/match/v5/matches/TW2_81882122/timeline"
	// buff, err := f.Get(url)

	// load local json data
	buff, err := os.ReadFile(path + "match_timeline.json")
	if err != nil {
		t.Fatal(err)
	}

	// parse
	var res *riotmodel.MatchTimelineDTO
	err = json.Unmarshal(buff, &res)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(unsafe.Sizeof(res))
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
		fmt.Printf("summoner %s's LP is %d\n", entries[i].SummonerID, entries[i].LeaguePoints)
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

func Test_parse_champion(t *testing.T) {
	lang := utils.ConvertLangToLangStr(riotmodel.LANG_zh_CN)
	version := "13.14.1"
	championName := "Aatrox"
	url := fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/data/%s/champion/%s.json",
		version, lang, championName)
	buff, err := f.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	var res *riotmodel.ChampionSingleDTO
	err = json.Unmarshal(buff, &res)
	if err != nil {
		t.Fatal(err)
	}
	champion := res.Data[championName]
	if err := db.Create(champion).Error; err != nil {
		t.Log(err)
	}
}

func Test_parse_champions(t *testing.T) {
	lang := utils.ConvertLangToLangStr(riotmodel.LANG_zh_CN)
	version := "13.14.1"
	url := fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/%s/data/%s/champion.json", version, lang)
	buff, err := f.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	var res *riotmodel.ChampionListDTO
	err = json.Unmarshal(buff, &res)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_parse_version(t *testing.T) {
	url := "https://ddragon.leagueoflegends.com/api/versions.json"
	buff, err := f.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	var res []string
	err = json.Unmarshal(buff, &res)
	if err != nil {
		t.Fatal(err)
	}
	curVersion := res[0]
	t.Log(curVersion)
}

func Test_parse_item_list(t *testing.T) {
	lang := utils.ConvertLangToLangStr(riotmodel.LANG_zh_CN)
	url := fmt.Sprintf("https://ddragon.leagueoflegends.com/cdn/13.14.1/data/%s/item.json", lang)
	buff, err := f.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	var res *riotmodel.ItemList
	err = json.Unmarshal(buff, &res)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_parse_spell_list(t *testing.T) {
	url := "https://ddragon.leagueoflegends.com/cdn/14.1.1/data/zh_CN/summoner.json"
	buff, err := f.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	var res *riotmodel.SpellList
	err = json.Unmarshal(buff, &res)
	if err != nil {
		t.Fatal(err)
	}
}
