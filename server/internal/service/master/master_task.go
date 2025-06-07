package master

import (
	"context"
	"time"

	"github.com/cralack/ChaosMetrics/server/pkg/xamqp"
	"github.com/cralack/ChaosMetrics/server/proto/publisher"
	"go-micro.dev/v4/client"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TaskSpec struct {
	ID           string
	Name         string
	AssignedNode string
	CreationTime int64

	SumName string
	Type    string
	Loc     string
}

func (m *Master) AddTask(assigner Assigner, tasks ...*TaskSpec) {
	for _, task := range tasks {
		if task.ID == "" {
			task.ID = m.IDGen.Generate().String()
		}

		if node, err := assigner.Assign(m.workNodes); err != nil {
			m.logger.Error("assign failed", zap.Error(err))
			continue
		} else {
			task.AssignedNode = node.Id
			if task.CreationTime == 0 {
				task.CreationTime = time.Now().Unix()
			}
			m.logger.Debug("add task", zap.Any("specs", task))
		}

		// producer here
		body := Encode(task)
		if err := m.producer.Publish([]byte(body), xamqp.Exchange, task.AssignedNode, 0); err != nil {
			m.logger.Error("publish task failed", zap.Error(err))
		} else {
			m.logger.Debug("publish task", zap.Any("specs", task))
		}
	}
}

var _ publisher.PublisherHandler = &Master{}

// PushTask implement grpc call
func (m *Master) PushTask(ctx context.Context, ptask *publisher.TaskSpec, out *emptypb.Empty) error {
	// mark 'out' as unused
	_ = out

	// for follower
	if !m.IsLeader() && m.leaderID != "" && m.leaderID != m.ID {
		addr := getLeaderAddr(m.leaderID)
		_, err := m.forwardCli.PushTask(ctx, ptask, client.WithAddress(addr))
		if err != nil {
			m.logger.Error("forward failed", zap.Error(err))
			return err
		}
		return nil
	}

	// for master
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
