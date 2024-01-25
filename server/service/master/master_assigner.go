package master

import (
	"errors"

	"go-micro.dev/v4/registry"
)

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

type AreaAssigner struct{}
