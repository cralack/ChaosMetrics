package test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

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
		cmd = append(cmd, pipe.HSet(context.Background(), "/entry/tw2", e.SummonerID, e))
	}
	if _, err := pipe.Exec(context.Background()); err != nil {
		t.Log(err)
	}
}
