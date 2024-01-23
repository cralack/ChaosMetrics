package init

import (
	"github.com/cralack/ChaosMetrics/server/internal/config"
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/internal/pkg/xgorm"
	"github.com/cralack/ChaosMetrics/server/internal/pkg/xredis"
	"github.com/cralack/ChaosMetrics/server/internal/pkg/xviper"
	"github.com/cralack/ChaosMetrics/server/internal/pkg/xzap"
	"github.com/cralack/ChaosMetrics/server/model/anres"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"go.uber.org/zap"
)

var err error

func init() {
	// set runtime envs
	global.GvaConf = config.New()

	// setup config service
	global.GvaVp, err = xviper.Viper()
	if err != nil {
		panic(err)
	}

	// setup logger service
	global.GvaLog, err = xzap.Zap(global.GvaEnv)
	if err != nil {
		panic(err)
	}

	// setup orm service
	global.GvaDb, err = xgorm.GetDB()
	if err != nil {
		panic(err)
	}

	// setup redis service
	global.GvaRdb, err = xredis.GetClient()
	if err != nil {
		panic(err)
	}

	global.GvaLog.Debug("env pkg init succeed")

	// if.need.AutoMigrate
	if err := global.GvaDb.AutoMigrate(
		&riotmodel.LeagueEntryDTO{},
		&riotmodel.SummonerDTO{},
		// match
		&riotmodel.MatchDB{},
		&riotmodel.ParticipantDB{},

		&anres.Champion{},
	); err != nil {
		global.GvaLog.Error("init orm model failed", zap.Error(err))
	} else {
		global.GvaLog.Debug("init orm model succeed")
	}
}
