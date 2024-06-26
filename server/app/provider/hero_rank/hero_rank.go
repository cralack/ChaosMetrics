package hero_rank

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

type HeroRankService struct {
	rdb   *redis.Client
	rlock *sync.RWMutex
}

func NewHeroRankService() *HeroRankService {
	return &HeroRankService{
		rdb:   global.ChaRDB,
		rlock: &sync.RWMutex{},
	}
}

func (s *HeroRankService) QueryHeroRank(version, location, gamemode string) ([]*anres.ChampionBrief, error) {
	var (
		res  []*anres.ChampionBrief
		err  error
		buff string
	)
	key := "/champion_brief"
	vidx, _ := utils.ConvertVersionToIdx(version)
	field := fmt.Sprintf("%d_%s@%s", vidx, gamemode, location)
	buff = s.rdb.HGet(context.Background(), key, field).Val()
	if err = json.Unmarshal([]byte(buff), &res); err != nil {
		return nil, err
	}
	return res, nil
}
