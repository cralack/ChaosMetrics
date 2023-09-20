package master

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/snowflake"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Master struct {
	ID        string
	Loc       string
	leaderID  string
	workNoeds map[string]int
	IDGen     *snowflake.Node
	etcdCli   *clientv3.Client
	rwlock    *sync.Mutex
}

func (*Master) Run() {
	fmt.Println("")
}
