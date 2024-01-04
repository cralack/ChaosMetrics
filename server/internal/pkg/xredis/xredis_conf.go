package xredis

import (
	"errors"
	"strconv"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/redis/go-redis/v9"
)

func GetRedisConf() (*redis.Options, error) {
	conf := global.GVA_CONF.RedisConf
	if conf.Host == "" || conf.Port == 0 {
		return nil, errors.New("get redis conf failed")
	}
	port := strconv.Itoa(conf.Port)
	option := &redis.Options{
		Addr:            conf.Host + ":" + port,
		DB:              conf.DB,
		Username:        conf.Username,
		Password:        conf.Password,
		DialTimeout:     conf.Timeout,
		ReadTimeout:     conf.ReadTimeout,
		WriteTimeout:    conf.WriteTimeout,
		MinIdleConns:    conf.ConnMinIdle,
		PoolSize:        conf.ConnMaxOpen,
		ConnMaxLifetime: conf.ConnMaxLifetime,
		ConnMaxIdleTime: conf.ConnMaxIdleTime,
	}
	return option, nil
}
