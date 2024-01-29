package init

import (
	"github.com/cralack/ChaosMetrics/server/internal/config"
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/pkg/xgorm"
	"github.com/cralack/ChaosMetrics/server/pkg/xredis"
	"github.com/cralack/ChaosMetrics/server/pkg/xviper"
	"github.com/cralack/ChaosMetrics/server/pkg/xzap"
	"go.uber.org/zap"
)

var err error

func init() {
	// set runtime envs
	global.ChaConf = config.New()

	// setup config service
	global.ChaViper, err = xviper.Viper()
	if err != nil {
		panic(err)
	}

	// setup logger service
	global.ChaLogger, err = xzap.Zap(global.ChaEnv)
	if err != nil {
		panic(err)
	}

	// setup orm service
	global.ChaDB, err = xgorm.GetDB()
	if err != nil {
		panic(err)
	}

	// setup redis service
	global.ChaRDB, err = xredis.GetClient()
	if err != nil {
		panic(err)
	}

	global.ChaLogger.Debug("env pkg init succeed")

	// if.need.AutoMigrate
	if err = global.ChaDB.AutoMigrate(
		&riotmodel.LeagueEntryDTO{},
		&riotmodel.SummonerDTO{},
		// match
		&riotmodel.MatchDB{},
		&riotmodel.ParticipantDB{},

		// &anres.ChampionDetail{},
	); err != nil {
		global.ChaLogger.Error("init orm model failed", zap.Error(err))
	} else {
		global.ChaLogger.Debug("init orm model succeed")
	}
}
