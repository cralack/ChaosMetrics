package master

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/cralack/ChaosMetrics/server/utils"
	clientv3 "go.etcd.io/etcd/client/v3"

	"go-micro.dev/v4/registry"
)

const (
	ServiceName = "pumper.worker"
)

type Master struct {
	ready     int32
	ID        string
	Loc       string
	leaderID  string
	workNodes map[string]*NodeSpec
	resources map[string]*ResourceSpec
	rlock     *sync.Mutex
	IDGen     *snowflake.Node
	etcdCli   *clientv3.Client

	options
}

type ResourceSpec struct {
	ID           string
	Name         string
	Loc          string
	AssignedNode string
	CreationTime int64
}

type NodeSpec struct {
	Node *registry.Node
}

func New(id string, opts ...Option) (*Master, error) {
	m := &Master{
		workNodes: make(map[string]*NodeSpec),
		resources: make(map[string]*ResourceSpec),
		rlock:     &sync.Mutex{},
	}
	options := defaultOptions
	for _, opt := range opts {
		opt(&options)
	}
	m.options = options

	ipv4, err := utils.GetLocalIP()
	if err != nil {
		return nil, err
	}
	m.ID = fmt.Sprintf("master_%s@%s%s", id, ipv4, m.GRPCAddress)
	m.logger.Debug("master id:" + m.ID)

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
	// start core func
	go m.Campaign()
}
