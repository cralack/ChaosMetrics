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

	// if.need.AutoMigrate
	if err := global.GVA_DB.AutoMigrate(
		&riotmodel.LeagueEntryDTO{},
		&riotmodel.SummonerDTO{},
		// match
		&riotmodel.MatchDB{},
		&riotmodel.ParticipantDB{},

		&anres.Champion{},
	); err != nil {
		global.GVA_LOG.Error("init orm model failed", zap.Error(err))
	} else {
		global.GVA_LOG.Debug("init orm model succeed")
	}
}
