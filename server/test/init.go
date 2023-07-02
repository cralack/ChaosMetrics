package test

import (
	"github.com/cralack/ChaosMetrics/server/global"
	_ "github.com/cralack/ChaosMetrics/server/init"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/pkg/fetcher"
	"go.uber.org/zap"
)

var f fetcher.Fetcher
var path string

func init() {
	f = fetcher.NewBrowserFetcher()
	path = "./local_json/"
	global.GVA_DB.Exec("DROP TABLE IF EXISTS match_dtos, match_summoners, participant_dtos, summoner_dtos, team_dtos")
	// AutoMigrate
	if err := global.GVA_DB.AutoMigrate(
		&riotmodel.MatchDto{},
		&riotmodel.ParticipantDto{},
		&riotmodel.TeamDto{},
		&riotmodel.SummonerDTO{},
		&riotmodel.LeagueItemDTO{},
	); err != nil {
		global.GVA_LOG.Error("init db model failed", zap.Error(err))
	}
}
