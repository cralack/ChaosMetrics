package champion_rank

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model/anres"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/redis/go-redis/v9"
)

type ChampionRankService struct {
	rdb   *redis.Client
	rlock *sync.RWMutex
}

func NewChampionRankService() *ChampionRankService {
	return &ChampionRankService{
		rdb:   global.ChaRDB,
		rlock: &sync.RWMutex{},
	}
}

func (c *ChampionRankService) QueryChampionRank(version, location, gamemode string) ([]*anres.ChampionBrief, error) {
	var (
		res  []*anres.ChampionBrief
		err  error
		buff string
	)
	key := "/champion_brief"
	vidx, _ := utils.ConvertVersionToIdx(version)
	field := fmt.Sprintf("%d_%s@%s", vidx, gamemode, location)
	buff = c.rdb.HGet(context.Background(), key, field).Val()
	if err = json.Unmarshal([]byte(buff), &res); err != nil {
		return nil, err
	}
	return res, nil
}
