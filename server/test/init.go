package test

import (
	"github.com/cralack/ChaosMetrics/server/global"
	_ "github.com/cralack/ChaosMetrics/server/init"
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
	// global.GVA_DB.Exec("DROP TABLE IF EXISTS match_dtos, match_summoners, participant_dtos, summoner_dtos, team_dtos, league_entry_dtos")
	// global.GVA_RDB.FlushDB(context.Background())
	// AutoMigrate
	if err := db.AutoMigrate(
		&riotmodel.LeagueEntryDTO{},
		&riotmodel.SummonerDTO{},
		&riotmodel.MatchDto{},
		&riotmodel.ParticipantDto{},
		&riotmodel.TeamDto{},
	); err != nil {
		logger.Error("init gormdb model failed", zap.Error(err))
	} else {
		logger.Info("init gormdb model succeed")
	}
}
