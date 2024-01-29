package champion_detail

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/internal/service/analyzer"
	"github.com/cralack/ChaosMetrics/server/model/anres"
	"github.com/redis/go-redis/v9"
)

type ChampionDetailService struct {
	rdb   *redis.Client
	rlock *sync.RWMutex
}

func NewChampionDetailService() *ChampionDetailService {
	return &ChampionDetailService{
		rdb:   global.ChaRDB,
		rlock: &sync.RWMutex{},
	}
}

func (c *ChampionDetailService) QueryChampionDetail(name, version, location, gamemode string) (*anres.ChampionDetail, error) {
	var (
		res  *anres.ChampionDetail
		err  error
		buff string
	)
	key := "/champion_detail"
	field := analyzer.GetID(name, version, location, gamemode)
	buff = c.rdb.HGet(context.Background(), key, field).Val()
	if err = json.Unmarshal([]byte(buff), &res); err != nil {
		return nil, err
	}

	return res, nil
}
