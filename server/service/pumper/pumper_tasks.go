package pumper

import (
	"context"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/service/master"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

func (p *Pumper) getTask() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // 设置超时
	defer cancel()

	resp, err := p.etcdcli.Get(
		ctx,
		global.TaskPath,
		clientv3.WithPrefix(),
	)
	if err != nil {
		p.logger.Error("etcd get task failed", zap.Error(err))
		return
	}

	for _, kv := range resp.Kvs {
		// unmarshal task
		t, err := master.Decode(kv.Value)
		if err != nil || t == nil {
			p.logger.Error("decode task failed", zap.Error(err))
			continue
		}

		// pop task to que
		id, err := master.GetNodeID(t.AssgnedNode)
		if err == nil && p.id == id {
			// push task to pumper que

			// delete task
			key := global.TaskPath + "/" + t.Name
			if resp2, err2 := p.etcdcli.Delete(ctx, key); err2 != nil || resp2.Deleted == 0 {
				p.logger.Error("pop task out failed", zap.Error(err2))
			} else {
				p.logger.Debug("delete succeed", zap.Any("resp", resp2))
			}
		}

	}
}

func (p *Pumper) watchTasks() {
	watcher := p.etcdcli.Watch(
		context.Background(),
		global.TaskPath,
		clientv3.WithPrefix(),
	)

	for w := range watcher {
		if w.Err() != nil {
			p.logger.Error("watch task failed", zap.Error(w.Err()))
			continue
		}
		if w.Canceled {
			p.logger.Error("watch task canceled")
			return
		}

		for _, event := range w.Events {
			switch event.Type {
			case clientv3.EventTypePut:
				// pop task
			}
		}
	}
}
