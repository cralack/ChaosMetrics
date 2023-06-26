package db

import (
	"ChaosMetrics/server/global"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetDB() (*gorm.DB, error) {
	if global.GVA_DB != nil {
		return global.GVA_DB, nil
	}
	//check dsn
	if global.GVA_CONF.Dbconf.DSN == "" {
		GetDBConfig()
	}
	//init val
	var (
		db       *gorm.DB
		gormConf *gorm.Config
		err      error
	)
	//diff logger for gom
	if global.GVA_ENV == global.PRODUCT_ENV {
		gormConf = &gorm.Config{
			Logger: &ZapLogger{
				global.GVA_LOG,
			},
		}
	} else {
		gormConf = &gorm.Config{}
	}
	//get db con
	db, err = gorm.Open(
		mysql.Open(global.GVA_CONF.Dbconf.DSN),
		gormConf,
	)
	if err != nil || db == nil {
		return nil, err
	}

	return db, nil
}
