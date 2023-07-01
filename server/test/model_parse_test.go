package test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
)

func Test_parse_summoner(t *testing.T) {
	// fetching remote JSON data (3~5 seconds per request)
	// url := "https://tw2.api.riotgames.com/lol/summoner/v4/summoners/by-name/Mudife"
	// buff, err := f.Get(url)
	buff, err := os.ReadFile(path + "summoner.txt")
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
	var res riotmodel.MatchDto
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
}
