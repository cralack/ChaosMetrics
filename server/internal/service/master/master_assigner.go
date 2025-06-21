package master

import (
	"errors"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"go-micro.dev/v5/registry"
)

type Assigner interface {
	Assign(workNodes map[string]*registry.Node) (*registry.Node, error)
}

type SimpleAssigner struct{}

var _ Assigner = &SimpleAssigner{}

// Assign :a random balacner
func (s *SimpleAssigner) Assign(workNodes map[string]*registry.Node) (*registry.Node, error) {
	// 检查是否有可用的工作节点
	size := len(workNodes)
	if size == 0 {
		return nil, errors.New("no worker available")
	}

	// 确保随机性

	idx := global.ChaRNG.Intn(size)

	// 获取第 idx 个节点（避免重复遍历）
	for _, n := range workNodes {
		if idx == 0 {
			return n, nil
		}
		idx--
	}

	// 正常情况下不应该到达这里
	return nil, errors.New("unexpected error: no worker found")
}

type AreaAssigner struct{}
