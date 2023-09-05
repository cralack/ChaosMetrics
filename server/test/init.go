package test

import (
	"github.com/cralack/ChaosMetrics/server/global"
	_ "github.com/cralack/ChaosMetrics/server/init"
	"github.com/cralack/ChaosMetrics/server/model/anres"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
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
	db = global.GVA_DB
	rdb = global.GVA_RDB
	logger = global.GVA_LOG
	// wipe gdb && rdb
	global.GVA_DB.Exec("DROP TABLE IF EXISTS analyzed_champions")
	// global.GVA_RDB.FlushDB(context.Background())
	// AutoMigrate
	if err := db.AutoMigrate(
		&riotmodel.LeagueEntryDTO{},
		&riotmodel.SummonerDTO{},
		// match
		&riotmodel.MatchDB{},
		&riotmodel.ParticipantDB{},

		&anres.Champion{},
	); err != nil {
		logger.Error("init xgorm model failed", zap.Error(err))
	} else {
		logger.Debug("init xgorm model succeed")
	}
}
