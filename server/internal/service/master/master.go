package master

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/pkg/xamqp"
	"github.com/cralack/ChaosMetrics/server/proto/publisher"
	"github.com/cralack/ChaosMetrics/server/utils"
	"go-micro.dev/v4/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Master struct {
	ready      int32
	ID         string
	Loc        string
	leaderID   string
	workNodes  map[string]*registry.Node
	rlock      *sync.RWMutex
	IDGen      *snowflake.Node
	etcdCli    *clientv3.Client
	producer   xamqp.MessageQueue
	forwardCli publisher.PublisherService

	options
}

func New(id string, opts ...Setup) (*Master, error) {
	m := &Master{
		workNodes: make(map[string]*registry.Node),
		rlock:     &sync.RWMutex{},
	}
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	m.options = options

	// init masterID

	if ipv4, err := utils.GetLocalIP(); err != nil {
		return nil, err
	} else {
		m.ID = fmt.Sprintf("%s@%s%s", id, ipv4, m.GRPCAddress)
		m.logger.Debug("master id:" + m.ID)
	}

	// init snowflake
	machineId, _ := strconv.Atoi(global.ChaConf.System.ID)
	if node, err := snowflake.NewNode(int64(machineId)); err != nil {
		return nil, err
	} else {
		m.IDGen = node
	}
	// init producer for master
	producer, err := xamqp.NewRabbitMQ(
		xamqp.Producer,
		nil,
		xamqp.WithContext(m.ctx),
	)
	if err != nil {
		return nil, errors.New("init producer for master fail:" + err.Error())
	}
	if err = producer.Start(); err != nil {
		return nil, errors.New("start producer for master fail:" + err.Error())
	}
	m.producer = producer

	// init master's etcd client
	if cli, err2 := clientv3.New(clientv3.Config{
		Endpoints: []string{m.registryURL},
	}); err2 != nil {
		return nil, err2
	} else {
		m.etcdCli = cli
	}
	m.updateWorkNodes()

	return m, nil
}

func (m *Master) Run() {
	// start elect func
	go m.Campaign()
}

func (m *Master) SetForwardCli(cli publisher.PublisherService) {
	m.forwardCli = cli
}
