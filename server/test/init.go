package test

import (
	"github.com/cralack/ChaosMetrics/server/global"
	_ "github.com/cralack/ChaosMetrics/server/init"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/service/fetcher"
	"go.uber.org/zap"
)

var f fetcher.Fetcher
var path string

func init() {
	f = fetcher.NewBrowserFetcher()
	path = "./local_json/"
	// wipe gdb && rdb
	// global.GVA_DB.Exec("DROP TABLE IF EXISTS match_dtos, match_summoners, participant_dtos, summoner_dtos, team_dtos, league_entry_dtos")
	// global.GVA_RDB.FlushDB(context.Background())
	// AutoMigrate
	if err := global.GVA_DB.AutoMigrate(
		&riotmodel.LeagueEntryDTO{},
		&riotmodel.SummonerDTO{},
		&riotmodel.MatchDto{},
		&riotmodel.ParticipantDto{},
		&riotmodel.TeamDto{},
	); err != nil {
		global.GVA_LOG.Error("init gormdb model failed", zap.Error(err))
	} else {
		global.GVA_LOG.Info("init gormdb model succeed")
	}
}
