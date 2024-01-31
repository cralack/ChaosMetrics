package item

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ItmService struct {
	rdb   *redis.Client
	rlock *sync.RWMutex
}

func NewItemService(params ...interface{}) *ItmService {
	if len(params) > 0 {
		global.ChaLogger.Debug("param", zap.Int("len", len(params)))
	}
	return &ItmService{
		rdb:   global.ChaRDB,
		rlock: &sync.RWMutex{},
	}
}

func (s *ItmService) QueryItem(itemId, version, lang string) (*riotmodel.ItemDTO, error) {
	var (
		item *riotmodel.ItemDTO
		vIdx uint
		err  error
		buff string
	)
	vIdx, err = utils.ConvertVersionToIdx(version)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("/items/%s", lang)
	field := fmt.Sprintf("%s@%d", itemId, vIdx)
	buff = s.rdb.HGet(context.Background(), key, field).Val()
	if err = json.Unmarshal([]byte(buff), &item); err != nil {
		return nil, err
	}
	return item, nil
}
