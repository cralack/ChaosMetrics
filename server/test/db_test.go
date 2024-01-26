package test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	gorm.Model
	Name string `gorm:"primary_key;column:user_name;type:varchar(100)"`
	Sex  bool
	Age  int
}

func Test_db_crud_func(t *testing.T) {
	// init xgorm
	db.Exec("DROP TABLE IF EXISTS users")
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatal(err)
	}

	// 增
	db.Create(&User{
		Name: "snoop",
		Sex:  false,
		Age:  18,
	})
	db.Create(&User{
		Name: "grant",
		Sex:  true,
		Age:  20,
	})
	db.Create(&User{
		Name: "rui",
		Sex:  true,
		Age:  22,
	})

	// 查
	var tar User
	db.First(&tar, "user_name=?", "grant")
	fmt.Println("first:\n", tar.Name, tar.ID)

	var tars []User
	// snop+grant		+rui
	db.Where("age<?", 21).Or("id>?", 2).Find(&tars)
	fmt.Println("find:")
	for _, t := range tars {
		fmt.Println(t.ID, t.Name, t.Age)
	}

	// 改
	db.Where("id=?", 1).First(&User{}).Update("user_name", "snoop dogg")
	db.Where("id in (?)", []int{1, 2}).Find(&[]User{}).Updates(
		map[string]interface{}{
			"Name": "after",
			"Sex":  true,
			"Age":  19,
		})

	// 删
	db.Where("id in (?)", []int{1, 3}).Delete(&User{})
	db.Where("id=?", 2).Unscoped().Delete(&User{})
}

func Test_match_store(t *testing.T) {
	// load local json data
	buff, err := os.ReadFile(path + "match.json")
	if err != nil {
		t.Fatal(err)
	}
	// parse
	var match *riotmodel.MatchDTO
	err = json.Unmarshal(buff, &match)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_purge_weak_player(t *testing.T) {
	level := "EMERALD"
	var entries []*riotmodel.LeagueEntryDTO
	if err := db.Find(&entries).Where("tier = ?", level).Error; err != nil {
		t.Log(err)
	}
	if err := db.Unscoped().Select(clause.Associations).Delete(entries).Error; err != nil {
		logger.Error("orm hard delete match failed ")
	}
	t.Log("succeed")
}

func Test_summoners_store(t *testing.T) {
	// load json
	buff, err := os.ReadFile(path + "summoners.json")
	if err != nil {
		t.Fatal(err)
	}

	// parse to model
	var summoners []*riotmodel.SummonerDTO
	err = json.Unmarshal(buff, &summoners)
	if err != nil {
		fmt.Println("解析失败：", err)
		return
	}
	// store
	db.Save(summoners)
}

func Test_summoner_entry_store(t *testing.T) {
	buff, err := os.ReadFile(path + "challenger_league.json")
	if err != nil {
		log.Fatal(err)
	}
	// parse to model
	var league *riotmodel.LeagueListDTO
	err = json.Unmarshal(buff, &league)
	if err != nil {
		t.Log("解析失败：", err)
	}
	var entry *riotmodel.LeagueEntryDTO
	for _, e := range league.Entries {
		if e.SummonerID == "CA6s2lmln5pY47sa0BcpyV39Vq7vDVSTMUVtJ17h7J3UURCTcDFwFS4lrg" {
			entry = e
		}
	}

	buff, err = os.ReadFile(path + "summoner.json")
	if err != nil {
		log.Fatal(err)
	}
	var summoner *riotmodel.SummonerDTO
	err = json.Unmarshal(buff, &summoner)
	if err != nil {
		t.Log("解析失败：", err)
	}
	db.Create(entry)

	var tar *riotmodel.LeagueEntryDTO
	db.Preload(clause.Associations).First(&tar)
	t.Log(tar.SummonerID)
}

func Test_isExist(t *testing.T) {
	ctx := context.Background()
	key := "/summoner/tw2"
	redisMap := make(map[string]*riotmodel.SummonerDTO)
	if size := rdb.HLen(ctx, key).Val(); size != 0 {
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

	entryExistMap := make(map[string]bool)
	keys := make([]string, 0, 2*len(redisMap))
	for k := range redisMap {
		keys = append(keys, k)
	}
	var existintEntry []*riotmodel.LeagueEntryDTO
	if err := db.Where("summoner_id IN ?", keys).Find(&existintEntry).Error; err != nil {
		t.Log(err)
	}
	for _, e := range existintEntry {
		entryExistMap[e.SummonerID] = true
	}
	for _, k := range keys {
		if _, has := entryExistMap[k]; !has {
			entryExistMap[k] = false
		}
	}
	logger.Debug("ok")

}

// may need setup gorm's logger silent before test
// server/pkg/xgorm/xgorm.go:38
// xgorm.Save([size]*riotmodel.LeagueEntryDTO) size=1~10k store benchmark
func Benchmark_db_store_1(b *testing.B) {
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
	// init test val
	type storeTEST struct {
		entires []*riotmodel.LeagueEntryDTO
		n       int
	}
	leagueSize := len(league.Entries)
	tests := make([]*storeTEST, 5)
	cnt := 1
	for i := range tests {
		size := cnt
		cnt *= 10
		tests[i] = &storeTEST{
			n:       size,
			entires: make([]*riotmodel.LeagueEntryDTO, size),
		}
		tt := tests[i]
		for idx := range tt.entires {
			tt.entires[idx] = league.Entries[i%leagueSize]
		}
	}
	// test
	db := global.GvaDb
	for _, tt := range tests {
		b.Run(fmt.Sprintf("%d entries", tt.n), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				db.Save(tt.entires)
			}
			b.StopTimer()
		})
	}

	/*
										local result
		goos: windows
		goarch: amd64
		pkg: github.com/cralack/ChaosMetrics/server/test
		cpu: AMD Ryzen 5 2600X Six-Core Processor
		Benchmark_db_store
		Benchmark_db_store/1_entries
		Benchmark_db_store/1_entries-12                      128          11071534 ns/op

		Benchmark_db_store/10_entries
		Benchmark_db_store/10_entries-12                     120          10154584 ns/op

		Benchmark_db_store/0.1k_entries
		Benchmark_db_store/0.1k_entries-12                    88          19098601 ns/op

		Benchmark_db_store/1k_entries
		Benchmark_db_store/1k_entries-12                      20          76719330 ns/op

		Benchmark_db_store/10k_entries
		Benchmark_db_store/10k_entries-12                     10         107688520 ns/op
	*/
}

// xgorm.Save([size]*riotmodel.LeagueEntryDTO) size=1k~10k store benchmark
func Benchmark_db_store_2(b *testing.B) {
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
	// init test val
	type storeTEST struct {
		entires []*riotmodel.LeagueEntryDTO
		n       int
	}
	leagueSize := len(league.Entries)
	tests := make([]*storeTEST, 10)
	for i := range tests {
		size := (1 + i) * 100
		tests[i] = &storeTEST{
			n:       size,
			entires: make([]*riotmodel.LeagueEntryDTO, size),
		}
		tt := tests[i]
		for idx := range tt.entires {
			tt.entires[idx] = league.Entries[i%leagueSize]
		}
	}
	// test
	db := global.GvaDb
	for _, tt := range tests {
		b.Run(fmt.Sprintf("%d entries", tt.n), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				db.Save(tt.entires)
			}
			b.StopTimer()
		})
	}

	/*
											local result
		goos: windows
		goarch: amd64
		pkg: github.com/cralack/ChaosMetrics/server/test
		cpu: AMD Ryzen 5 2600X Six-Core Processor
		Benchmark_db_store_2
		Benchmark_db_store_2/100_entries
		Benchmark_db_store_2/100_entries-12                   97          15569600 ns/op

		Benchmark_db_store_2/200_entries
		Benchmark_db_store_2/200_entries-12                   57          22760793 ns/op

		Benchmark_db_store_2/300_entries
		Benchmark_db_store_2/300_entries-12                   54          25559107 ns/op

		Benchmark_db_store_2/400_entries
		Benchmark_db_store_2/400_entries-12                   54          37813459 ns/op

		Benchmark_db_store_2/500_entries
		Benchmark_db_store_2/500_entries-12                   40          40469692 ns/op

		Benchmark_db_store_2/600_entries
		Benchmark_db_store_2/600_entries-12                   24          59813025 ns/op

		Benchmark_db_store_2/700_entries
		Benchmark_db_store_2/700_entries-12                   26          55035704 ns/op

		Benchmark_db_store_2/800_entries
		Benchmark_db_store_2/800_entries-12                   21          54065381 ns/op

		Benchmark_db_store_2/900_entries
		Benchmark_db_store_2/900_entries-12                   25          62689404 ns/op

		Benchmark_db_store_2/1000_entries
		Benchmark_db_store_2/1000_entries-12                  19          65841416 ns/op

	*/
}
