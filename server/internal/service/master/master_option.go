package master

import (
	"context"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"go-micro.dev/v5/registry"
	"go.uber.org/zap"
)

type options struct {
	logger      *zap.Logger
	registry    registry.Registry
	ctx         context.Context
	registryURL string
	GRPCAddress string
}

var defaultOptions = options{
	logger: global.ChaLogger,
}

type Setup func(opts *options)

func WithLogger(logger *zap.Logger) Setup {
	return func(opts *options) {
		opts.logger = logger
	}
}

func WithregistryURL(registryURL string) Setup {
	return func(opts *options) {
		opts.registryURL = registryURL
	}
}

func WithRegistry(registry registry.Registry) Setup {
	return func(opts *options) {
		opts.registry = registry
	}
}

func WithGRPCAddress(GRPCAddress string) Setup {
	return func(opts *options) {
		opts.GRPCAddress = GRPCAddress
	}
}

func WithContext(ctx context.Context) Setup {
	return func(opts *options) {
		opts.ctx = ctx
	}
}
