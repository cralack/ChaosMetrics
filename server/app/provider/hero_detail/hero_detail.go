package hero_detail

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/redis/go-redis/v9"
)

type HeroDetailService struct {
	rdb *redis.Client
}

func NewHeroDetailService() *HeroDetailService {
	return &HeroDetailService{
		rdb: global.ChaRDB,
	}
}

func (s *HeroDetailService) QueryHeroDetail(name, version, lang string) (*riotmodel.ChampionDTO, error) {
	var (
		res  *riotmodel.ChampionDTO
		err  error
		buff string
	)
	key := fmt.Sprintf("/champions/%s", lang)
	vidx, _ := utils.ConvertVersionToIdx(version)
	field := fmt.Sprintf("%s@%d", name, vidx)
	buff = s.rdb.HGet(context.Background(), key, field).Val()
	if err = json.Unmarshal([]byte(buff), &res); err != nil {
		return nil, err
	}
	return res, nil
}
