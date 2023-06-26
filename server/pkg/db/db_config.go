package db

import (
	"net"
	"strconv"
	"time"

	"ChaosMetrics/server/global"

	"github.com/go-sql-driver/mysql"
)

func GetDBConfig() error {
	dbConf := global.GVA_CONF.Dbconf
	port := strconv.Itoa(dbConf.Port)
	timeout, err := time.ParseDuration(dbConf.Timeout)
	if err != nil {
		return err
	}
	readTimeout, err := time.ParseDuration(dbConf.ReadTimeout)
	if err != nil {
		return err
	}
	writeTimeout, err := time.ParseDuration(dbConf.WriteTimeout)
	if err != nil {
		return err
	}
	location, err := time.LoadLocation(dbConf.Loc)
	if err != nil {
		return err
	}
	driverConf := mysql.Config{
		User:                 dbConf.Username,
		Passwd:               dbConf.Password,
		Net:                  dbConf.Protocol,
		Addr:                 net.JoinHostPort(dbConf.Host, port),
		DBName:               dbConf.DBName,
		Collation:            dbConf.Collation,
		Loc:                  location,
		Timeout:              timeout,
		ReadTimeout:          readTimeout,
		WriteTimeout:         writeTimeout,
		ParseTime:            dbConf.ParseTime,
		AllowNativePasswords: true,
	}
	dbConf.DSN = driverConf.FormatDSN()

	global.GVA_CONF.Dbconf = dbConf
	return nil
}
