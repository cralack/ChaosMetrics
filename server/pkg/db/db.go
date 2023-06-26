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

	if global.GVA_CONF.Dbconf.DSN == "" {
		GetDBConfig()
	}
	var db *gorm.DB
	var err error
	db, err = gorm.Open(mysql.Open(global.GVA_CONF.Dbconf.DSN), &gorm.Config{})
	if err != nil || db == nil {
		return nil, err
	}

	return db, nil
}
