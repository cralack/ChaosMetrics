package master

import (
	"context"
	"errors"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"go-micro.dev/v4/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type Command int

const (
	MSGADD Command = iota
	MSGDEL
)

type Message struct {
	Cmd Command
}

type TaskSpec struct {
	ID           string
	Name         string
	AssgnedNode  string
	CreationTime int64

	SumName string
	Type    string
	Loc     string
}

func (m *Master) AddTask(assigner Assigner, tasks ...*TaskSpec) {
	for _, task := range tasks {
		task.ID = m.IDGen.Generate().String()

		if node, err := assigner.Assign(m.workNodes); err != nil {
			m.logger.Error("assign failed", zap.Error(err))
			continue
		} else {
			task.AssgnedNode = node.Id + "|" + node.Address
			task.CreationTime = time.Now().Unix()
			m.logger.Debug("add task", zap.Any("specs", task))
		}

		if _, err := m.etcdCli.Put(context.Background(),
			getTaskPath(task.Name), Encode(task)); err != nil {
			m.logger.Error("put etcd failed", zap.Error(err))
			continue
		}
	}
}

func (m *Master) loadTask() error {
	resp, err := m.etcdCli.Get(context.Background(), global.TaskPath, clientv3.WithPrefix())
	if err != nil {
		return errors.New("master get task from etcd failed")
	}
	tasks := make(map[string]*TaskSpec)
	for _, kv := range resp.Kvs {
		if t, err2 := Decode(kv.Value); err2 == nil && t != nil {
			tasks[t.Name] = t
			timeStamp := time.Unix(
				t.CreationTime/int64(time.Second),
				t.CreationTime%int64(time.Second),
			)
			m.logger.Debug(timeStamp.String())
		}
	}

	m.rlock.Lock()
	defer m.rlock.Unlock()
	m.tasks = tasks

	m.logger.Info("leader init load task", zap.Int("lenth", len(m.tasks)))

	for _, t := range m.tasks {
		if t.AssgnedNode != "" {
			id, err := GetNodeID(t.AssgnedNode)
			if err != nil {
				m.logger.Error("get node ID faild", zap.Error(err))
			}
			if node, has := m.workNodes[id]; has {
				m.logger.Debug(node.Id)
			}
		}
	}
	return nil
}

type Assigner interface {
	Assign(workNodes map[string]*registry.Node) (*registry.Node, error)
}

type SimpleAssigner struct{}

var _ Assigner = &SimpleAssigner{}

func (s *SimpleAssigner) Assign(workNodes map[string]*registry.Node) (*registry.Node, error) {
	for _, n := range workNodes {
		return n, nil
	}
	return nil, errors.New("no worker available")
}

// assign by area
