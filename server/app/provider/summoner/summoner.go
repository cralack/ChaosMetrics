package summoner

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/internal/service/pumper"
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/proto/publisher"
	"github.com/redis/go-redis/v9"
	"go-micro.dev/v5"
	"go-micro.dev/v5/client"
	grpccli "go-micro.dev/v5/client/grpc"
	"go-micro.dev/v5/registry"
	"go-micro.dev/v5/registry/etcd"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SumonerService struct {
	db     *gorm.DB
	rdb    *redis.Client
	reg    registry.Registry
	cli    publisher.PublisherService
	logger *zap.Logger
	rlock  *sync.RWMutex
}

func NewSumnService(params ...interface{}) *SumonerService {
	if len(params) > 0 {
		global.ChaLogger.Debug("param", zap.Int("len", len(params)))
	}
	conf := global.ChaConf.Micro
	reg := etcd.NewEtcdRegistry(registry.Addrs(conf.RegistryAddress))
	name := "gin.grpc.client"

	server := micro.NewService(
		micro.Name(name),
		micro.Registry(reg),
		micro.Client(grpccli.NewClient()),
	)
	cli := publisher.NewPublisherService(name, server.Client())

	return &SumonerService{
		db:     global.ChaDB,
		rdb:    global.ChaRDB,
		reg:    reg,
		cli:    cli,
		logger: global.ChaLogger,
		rlock:  &sync.RWMutex{},
	}
}

func (s *SumonerService) QuerySummonerByName(puuid, loc string) *response.SummonerDTO {
	var (
		res      *riotmodel.SummonerDTO
		services []*registry.Service
		err      error
		buff     string
	)
	// query from redis
	key := fmt.Sprintf("/summoner/%s", loc)
	buff = s.rdb.HGet(context.Background(), key, puuid).Val()
	if err = json.Unmarshal([]byte(buff), &res); err == nil && res.PUUID == puuid {
		return s.HandleSummoner(res)
	}

	// query from db
	if err = s.db.Where("loc=?", loc).Where("name=?",
		puuid).Find(&res).Error; err == nil && res.PUUID == puuid {
		return s.HandleSummoner(res)
	}

	// reg
	services, err = s.reg.GetService(global.MasterServiceName)
	if err != nil || len(services) == 0 || len(services[0].Nodes) == 0 {
		s.logger.Error("no master node available", zap.Error(err))
		return nil
	}
	addr := services[0].Nodes[0].Address

	// fetch summoner
	_, e := s.cli.PushTask(context.Background(),
		&publisher.TaskSpec{
			Name:    "summoner_query_task",
			Loc:     loc,
			Sumname: puuid,
			Type:    pumper.SummonerTypeKey,
		},
		client.WithAddress(addr),
	)

	if e != nil {
		s.logger.Error("failed to push task", zap.Error(e))
	}

	return nil
}
