package xgorm

import (
	"log"
	"os"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"gorm.io/gorm/logger"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetDB() (*gorm.DB, error) {
	if global.GVA_DB != nil {
		return global.GVA_DB, nil
	}
	// check dsn
	if global.GVA_CONF.Dbconf.DSN == "" {
		err := GetDBConfig()
		if err != nil {
			global.GVA_LOG.Error("get xgorm config failed",
				zap.Error(err))
		}
	}
	// init val
	var (
		db       *gorm.DB
		gormConf *gorm.Config
		err      error
	)
	// diff logger for gom
	if global.GVA_ENV == global.PRODUCT_ENV {
		gormConf = &gorm.Config{
			Logger: &ZapLogger{
				global.GVA_LOG,
			},
		}
	} else {
		defaultLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             500 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		})

		gormConf = &gorm.Config{
			// logger: logger.Default.LogMode(logger.Silent), // 禁用日志输出
			Logger: defaultLogger,
		}
	}
	// get gorm con
	db, err = gorm.Open(
		mysql.Open(global.GVA_CONF.Dbconf.DSN),
		gormConf,
	)
	if err != nil || db == nil {
		return nil, err
	}

	return db, nil
}
