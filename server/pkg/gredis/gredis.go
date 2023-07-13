package gredis

import (
	"github.com/redis/go-redis/v9"
)

func GetClient() (*redis.Client, error) {
	option, err := GetRedisConf()
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(option)
	return client, nil
}
