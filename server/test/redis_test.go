package test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/cralack/ChaosMetrics/server/app/provider/summoner"
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/redis/go-redis/v9"
)

func Test_redis_crud(t *testing.T) {
	// load json
	buff, err := os.ReadFile(path + "challenger_league.json")
	if err != nil {
		log.Fatal(err)
	}

	// parse to model
	var league *riotmodel.LeagueListDTO
	err = json.Unmarshal(buff, &league)
	if err != nil {
		fmt.Println("解析失败：", err)
		return
	}
	entries := league.Entries
	// init rdb
	pipe := rdb.Pipeline()
	// create && update
	cmd := make([]*redis.IntCmd, 0, len(entries))
	for _, e := range entries {
		cmd = append(cmd, pipe.HSet(context.Background(), "/entry/tw2", e.Puuid, e))
	}
	if _, err = pipe.Exec(context.Background()); err != nil {
		t.Log(err)
	}
}

func Test_redis(t *testing.T) {
	buff := rdb.HGet(context.Background(), "/championlist", "1305").Val()
	cList := make([]string, 0)
	if err := json.Unmarshal([]byte(buff), &cList); err != nil {
		t.Log(err)
	}
	t.Log(cList[0])
}

func Test_query(t *testing.T) {
	s := summoner.NewSumnService()
	if sumn := s.QuerySummonerByName("na1", "xyaaaz"); sumn != nil {
		t.Log(sumn.Name)
	}
}

func Test_item_list(t *testing.T) {
	var (
		res []*riotmodel.ItemDTO
		err error
	)
	values := global.ChaRDB.HGet(context.Background(), "/items", "1401-aram-zh_CN").Val()
	if err = json.Unmarshal([]byte(values), &res); err != nil {
		t.Log(err)
	}
}
