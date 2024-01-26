package test

import (
	_ "github.com/cralack/ChaosMetrics/server/init"
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/service/fetcher"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	f      fetcher.Fetcher
	path   string
	db     *gorm.DB
	rdb    *redis.Client
	logger *zap.Logger
)

func init() {
	f = fetcher.NewBrowserFetcher()
	path = "./local_json/"
	db = global.ChaDB
	rdb = global.ChaRDB
	logger = global.ChaLogger
	// wipe gdb && rdb
	// global.ChaDB.Exec("DROP TABLE IF EXISTS analyzed_champions,entries,match_participants,matches,summoners")
	// global.ChaRDB.FlushDB(context.Background())
	// AutoMigrate
	// if err := db.AutoMigrate(
	// 	&riotmodel.LeagueEntryDTO{},
	// 	&riotmodel.SummonerDTO{},
	// 	// match
	// 	&riotmodel.MatchDB{},
	// 	&riotmodel.ParticipantDB{},
	//
	// 	&anres.Champion{},
	// ); err != nil {
	// 	logger.Error("init orm model failed", zap.Error(err))
	// } else {
	// 	logger.Debug("init orm model succeed")
	// }
}
