package init

import (
	"github.com/cralack/ChaosMetrics/server/config"
	"github.com/cralack/ChaosMetrics/server/global"
	"github.com/cralack/ChaosMetrics/server/pkg/xgorm"
	"github.com/cralack/ChaosMetrics/server/pkg/xredis"
	"github.com/cralack/ChaosMetrics/server/pkg/xviper"
	"github.com/cralack/ChaosMetrics/server/pkg/xzap"
)

var err error

func init() {
	// set runtime envs
	global.GVA_CONF = config.New()

	// setup config service
	global.GVA_VP, err = xviper.Viper()
	if err != nil {
		panic(err)
	}

	// setup logger service
	global.GVA_LOG, err = xzap.Zap(global.GVA_ENV)
	if err != nil {
		panic(err)
	}

	// setup orm service
	global.GVA_DB, err = xgorm.GetDB()
	if err != nil {
		panic(err)
	}

	// setup redis service
	global.GVA_RDB, err = xredis.GetClient()
	if err != nil {
		panic(err)
	}

	global.GVA_LOG.Debug("env pkg init succeed")

	// // if.need.AutoMigrate
	// if err := global.GVA_DB.AutoMigrate(
	// 	&riotmodel.MatchDto{},
	// 	&riotmodel.ParticipantDto{},
	// 	&riotmodel.TeamDto{},
	// 	&riotmodel.SummonerDTO{},
	// 	&riotmodel.LeagueEntryDTO{},
	// ); err != nil {
	// 	global.GVA_LOG.Error("init xgorm model failed", zap.Error(err))
	// } else {
	// 	global.GVA_LOG.Info("init xgorm model succeed")
	// }
}
