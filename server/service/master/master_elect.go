package master

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"go-micro.dev/v4/registry"

	"go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/zap"
)

func (m *Master) Campaign() {
	// create etcd client & session
	s, err := concurrency.NewSession(m.etcdCli, concurrency.WithTTL(5))
	if err != nil {
		m.logger.Error("start concurrency session failed", zap.Error(err))
	}
	defer func(s *concurrency.Session) {
		err2 := s.Close()
		if err2 != nil {
			m.logger.Error("stop concurrency session failed", zap.Error(err2))
		}
	}(s)

	e := concurrency.NewElection(s, global.ElectPath)
	leaderCh := make(chan error)

	// if leader exist,goroutine blocking
	// elseif current master elect succeed,leaderCh <- sig
	go m.elect(e, leaderCh)

	// check leader
	leaderChange := e.Observe(context.Background())
	select {
	case resp := <-leaderChange:
		m.logger.Info("leader change",
			zap.String("leader:", string(resp.Kvs[0].Value)))
	}
	workerNodeChange := m.WatchWorker()

	for {
		select {
		// handle election error
		case err2 := <-leaderCh:
			if err2 != nil {
				m.logger.Error("leader elect failed", zap.Error(err2))
				go m.elect(e, leaderCh)
			} else {
				m.logger.Info(m.ID + " start change to leader")
				m.leaderID = m.ID
				if !m.IsLeader() {
					if err3 := m.BecomeLeader(); err3 != nil {
						m.logger.Error("become leader failed", zap.Error(err3))
					}
				}
			}

		// handle leader changes
		case resp := <-leaderChange:
			if len(resp.Kvs) > 0 {
				m.logger.Info("watch leader change",
					zap.String("leader:", string(resp.Kvs[0].Value)))
				m.leaderID = string(resp.Kvs[0].Value)
				if m.ID != string(resp.Kvs[0].Value) {
					atomic.StoreInt32(&m.ready, 0)
				}
			}

		// handle worker changes
		case resp := <-workerNodeChange:
			m.logger.Info("watch worker change", zap.String("worker:",
				fmt.Sprintf("%s,action:%s", resp.Service.Name, resp.Action)))
			m.updateWorkNodes()

		case <-time.After(3 * time.Second):
			m.AddTask(&SimpleAssigner{}, &TaskSpec{
				Name:    "TEST",
				SumName: "Mes",
				Type:    "match",
				Loc:     "tw2",
			})

		// check leader every 30s
		case <-time.After(30 * time.Second):
			rsp, err2 := e.Leader(context.Background())
			if err2 != nil {
				m.logger.Info("get leader failed", zap.Error(err2))
				if errors.Is(err2, concurrency.ErrElectionNoLeader) {
					go m.elect(e, leaderCh)
				}
			}

			if rsp != nil && len(rsp.Kvs) > 0 {
				m.logger.Debug("get leader", zap.String("value", string(rsp.Kvs[0].Value)))
				if m.IsLeader() && m.ID != string(rsp.Kvs[0].Value) {
					atomic.StoreInt32(&m.ready, 0)
				}
			}
		}
	}
}

// blocking until elect succeed
func (m *Master) elect(e *concurrency.Election, ch chan error) {
	err := e.Campaign(context.Background(), m.ID)
	ch <- err
}

func (m *Master) IsLeader() bool {
	return atomic.LoadInt32(&m.ready) != 0
}

func (m *Master) BecomeLeader() error {
	m.updateWorkNodes()
	if err := m.loadTask(); err != nil {
		return err
	}
	atomic.StoreInt32(&m.ready, 1)
	return nil
}

func (m *Master) WatchWorker() chan *registry.Result {
	watcher, err := m.registry.Watch(registry.WatchService(global.WorkerServiceName))
	if err != nil {
		m.logger.Panic("watch worker failed", zap.Error(err))
	}
	ch := make(chan *registry.Result)
	go func() {
		for {
			res, err := watcher.Next()
			if err != nil {
				m.logger.Error("watch worker service failed", zap.Error(err))
				continue
			}
			ch <- res
		}
	}()
	return ch
}

func (m *Master) updateWorkNodes() {
	services, err := m.registry.GetService(global.WorkerServiceName)
	if err != nil {
		m.logger.Error("get worker list failed", zap.Error(err))
	}
	m.rlock.Lock()
	defer m.rlock.Unlock()

	nodes := make(map[string]*registry.Node)
	if len(services) > 0 {
		for _, spec := range services[0].Nodes {
			nodes[spec.Id] = spec
		}
	}

	added, deleted, changed := workNodeDiff(m.workNodes, nodes)
	m.logger.Sugar().Info(
		"worker joined: ", added,
		", leaved: ", deleted,
		", changed: ", changed)
	m.workNodes = nodes
}
