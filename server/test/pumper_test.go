package test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	
	"github.com/cralack/ChaosMetrics/server/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/service/pumper"
	"github.com/redis/go-redis/v9"
)

func Test_pumper_init_entry(t *testing.T) {
	puper := pumper.NewPumper(
		pumper.WithEndMark(riotmodel.DIAMOND, 1),
	)
	puper.InitEntries()
	Test_db_store_check(t)
}

func Test_db_store_check(t *testing.T) {
	var (
		count int64
		tier  uint
	)
	db := global.GVA_DB
	res := db.Model(&riotmodel.LeagueEntryDTO{}).Count(&count)
	if err := res.Error; err != nil {
		t.Log("count failed:", err)
	} else {
		t.Log("gorm db total store:", count)
	}
	
	// count
	for tier = riotmodel.CHALLENGER; tier <= riotmodel.IRON; tier++ {
		tierSTR, _ := pumper.ConvertRankToStr(tier, 1)
		var cnt int64
		res := global.GVA_DB.Model(&riotmodel.LeagueEntryDTO{}).Where("tier=?", tierSTR).Count(&cnt)
		if res.Error != nil {
			t.Log(res.Error)
		} else {
			t.Log(tierSTR, cnt)
		}
	}
}

func Test_cache_hash_store_check(t *testing.T) {
	rdb := global.GVA_RDB
	// load from redis
	key := "/entry/tw2"
	// rdb.HLen(), rdb.HVals()
	keys, err := rdb.HKeys(context.Background(), key).Result()
	if err != nil {
		t.Log(err)
	}
	t.Log(len(keys))
	entries := make(map[string]*riotmodel.LeagueEntryDTO, len(keys))
	cnt := 0
	res := rdb.HGetAll(context.Background(), key).Val()
	for k, r := range res {
		var tmp riotmodel.LeagueEntryDTO
		_ = json.Unmarshal([]byte(r), &tmp)
		entries[k] = &tmp
		if tmp.MiniSeries != nil {
			cnt++
		}
	}
	t.Log(cnt)
}

func Test_redis_crud(t *testing.T) {
	// load json
	buff, err := os.ReadFile(path + "challenger_league.txt")
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
	rdb := global.GVA_RDB
	pipe := rdb.Pipeline()
	cmd := make([]*redis.IntCmd, len(entries))
	for _, e := range entries {
		cmd = append(cmd, pipe.HSet(context.Background(), "/entry/tw2", e.SummonerID, e))
	}
	if _, err := pipe.Exec(context.Background()); err != nil {
		t.Log(err)
	}
}