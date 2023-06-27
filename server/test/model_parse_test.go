package test

import (
	"ChaosMetrics/server/model/riotmodel"
	"ChaosMetrics/server/pkg/fetcher"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

var f fetcher.Fetcher
var path string

func init() {
	f = fetcher.NewBrowserFetcher()
	path = "./local_json/"
}

func Test_parse_summoner(t *testing.T) {
	//fetching remote JSON data (3~5 seconds per request)
	// url := "https://tw2.api.riotgames.com/lol/summoner/v4/summoners/by-name/Mudife"
	// buff, err := f.Get(url)
	buff, err := ioutil.ReadFile(path + "summoner.txt")
	if err != nil {
		t.Fatal(err)
	}
	var res riotmodel.SummonerDTO
	err = json.Unmarshal(buff, &res)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("ID:", res.ID)
	fmt.Println("Account ID:", res.AccountID)
	fmt.Println("PUUID:", res.PUUID)
	fmt.Println("Name:", res.Name)
	fmt.Println("Profile Icon ID:", res.ProfileIconID)
	fmt.Println("Revision Date:", res.RevisionDate)
	fmt.Println("Summoner Level:", res.SummonerLevel)
}

func Test_parse_championr_rotation(t *testing.T) {
	//fetching remote JSON data (3~5 seconds per request)
	// url := "https://tw2.api.riotgames.com/lol/platform/v3/champion-rotations"
	// buff, err := f.Get(url)

	//load local json data
	buff, err := ioutil.ReadFile(path + "championr_rotation.txt")
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
	//fetching remote JSON data (3~5 seconds per request)
	// url := "https://tw2.api.riotgames.com/lol/champion-mastery/v4/champion-masteries/by-puuid/F4fFtqehQLBj8U5sKBZF--k-7akbtb1IX790lRd4whPI4pXDAuVyfswHetg2lz_kMe2NJ0gUo5EIig/top"
	// buff, err := f.Get(url)
	buff, err := ioutil.ReadFile(path + "championr_mastery.txt")
	if err != nil {
		t.Fatal(err)
	}
	var res []riotmodel.ChampionMasteryDto
	err = json.Unmarshal([]byte(buff), &res)
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
	//fetching remote JSON data (3~5 seconds per request)
	// url := "https://sea.api.riotgames.com/lol/match/v5/matches/TW2_81882122"
	// buff, err := f.Get(url)

	//load local json data
	buff, err := ioutil.ReadFile(path + "match.txt")
	if err != nil {
		t.Fatal(err)
	}
	//parse
	var res riotmodel.MatchDto
	err = json.Unmarshal(buff, &res)
	if err != nil {
		t.Fatal(err)
	}
	//check data
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

}
