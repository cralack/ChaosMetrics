package master

import (
	"context"
	"errors"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/proto/publisher"
	"github.com/golang/protobuf/ptypes/empty"
	"go-micro.dev/v4/client"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

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
		}
	}

	m.rlock.Lock()
	defer m.rlock.Unlock()
	m.tasks = tasks

	m.logger.Info("leader init load task", zap.Int("lenth", len(m.tasks)))

	for _, t := range m.tasks {
		m.AddTask(&SimpleAssigner{}, t)
	}
	return nil
}

var _ publisher.PublisherHandler = &Master{}

func (m *Master) PushTask(ctx context.Context, ptask *publisher.TaskSpec, out *empty.Empty) error {
	// mark 'out' as unused
	_ = out

	if !m.IsLeader() && m.leaderID != "" && m.leaderID != m.ID {
		addr := getLeaderAddr(m.leaderID)
		_, err := m.forwardCli.PushTask(ctx, ptask, client.WithAddress(addr))
		m.logger.Error("forward failed", zap.Error(err))
		return err
	}
	m.rlock.Lock()
	defer m.rlock.Unlock()

	task := &TaskSpec{
		Name:    ptask.Name,
		SumName: ptask.Sumname,
		Type:    ptask.Type,
		Loc:     ptask.Loc,
	}
	m.AddTask(&SimpleAssigner{}, task)
	return nil
}
