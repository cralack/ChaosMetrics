package master

import (
	"github.com/bwmarrin/snowflake"
	"github.com/cralack/ChaosMetrics/server/utils"
	"go-micro.dev/v4/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Master struct {
	ID       string
	ready    int32
	leaderID string
	IDGen    *snowflake.Node
	etcdCli  *clientv3.Client
}

type ResourceSpec struct {
	ID           string
	Name         string
	Loc          string
	AssignedNode string
	CreationTime int64
}

type NodeSpec struct {
	Node    *registry.Node
	Address int
}

func New(id string) (*Master, error) {
	m := &Master{}

	ipv4, err := utils.GetLocalIP()
	if err != nil {
		return nil, err
	}
	m.ID = id + ipv4

	return m, nil
}
