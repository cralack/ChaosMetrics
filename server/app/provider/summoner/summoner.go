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
	grpccli "github.com/go-micro/plugins/v4/client/grpc"
	"github.com/go-micro/plugins/v4/registry/etcd"
	"github.com/redis/go-redis/v9"
	"go-micro.dev/v4"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/registry"
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
	conf := global.ChaConf.System
	reg := etcd.NewRegistry(registry.Addrs(conf.RegistryAddress))
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

func (s *SumonerService) QuerySummonerByName(name, loc string) *response.SummonerDTO {
	var (
		res      *riotmodel.SummonerDTO
		services []*registry.Service
		err      error
		buff     string
	)
	// query from redis
	key := fmt.Sprintf("/summoner/%s", loc)
	buff = s.rdb.HGet(context.Background(), key, name).Val()
	if err = json.Unmarshal([]byte(buff), &res); err == nil && res.Name == name {
		return s.HandleSummoner(res)
	}

	// query from db
	if err = s.db.Where("loc=?", loc).Where("name=?",
		name).Find(&res).Error; err == nil && res.Name == name {
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
			Sumname: name,
			Type:    pumper.SummonerTypeKey,
		},
		client.WithAddress(addr),
	)

	if e != nil {
		s.logger.Error("failed to push task", zap.Error(e))
	}

	return nil
}
