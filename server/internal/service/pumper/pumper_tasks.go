package pumper

import (
	"context"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/internal/service/master"
	"github.com/cralack/ChaosMetrics/server/utils"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

func (p *Pumper) getTask() {
	ctx := context.Background()

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
		task, err2 := master.Decode(kv.Value)
		if err2 != nil || task == nil {
			p.logger.Error("decode task failed", zap.Error(err2))
			continue
		}
		// pop task to pumper que
		p.handleTask(ctx, task)
	}
}

func (p *Pumper) watchTasks() {
	ctx := context.Background()
	watcher := p.etcdcli.Watch(
		ctx,
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
				task, err := master.Decode(event.Kv.Value)
				if err != nil || task == nil {
					p.logger.Error("decode task failed", zap.Error(err))
					continue
				}
				p.handleTask(ctx, task)
			}
		}
	}
}

func (p *Pumper) handleTask(ctx context.Context, task *master.TaskSpec) {
	var (
		err error
		id  string
	)
	id, err = master.GetNodeID(task.AssgnedNode)
	if err != nil || p.id != id {
		return
	}
	loc := utils.ConverHostLoCode(task.Loc)
	switch task.Type {
	case entryTypeKey:
		err = p.FetchEntryByName(task.SumName, loc)
	case matchTypeKey:
		err = p.FetchMatchByName(task.SumName, loc)
	}

	if err == nil {
		// such as: /tasks/TEST
		key := global.TaskPath + "/" + task.Name
		if resp2, err2 := p.etcdcli.Delete(ctx, key); err2 != nil || resp2.Deleted == 0 {
			p.logger.Error("pop task out failed", zap.Error(err2))
		} else {
			p.logger.Debug("delete succeed", zap.Any("resp", ""))
		}
	}
}
