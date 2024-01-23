package master

import (
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"go-micro.dev/v4/registry"
	"go.uber.org/zap"
)

type options struct {
	logger      *zap.Logger
	registry    registry.Registry
	registryURL string
	GRPCAddress string
}

var defaultOptions = options{
	logger: global.GvaLog,
}

type Option func(opts *options)

func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func WithregistryURL(registryURL string) Option {
	return func(opts *options) {
		opts.registryURL = registryURL
	}
}

func WithRegistry(registry registry.Registry) Option {
	return func(opts *options) {
		opts.registry = registry
	}
}

func WithGRPCAddress(GRPCAddress string) Option {
	return func(opts *options) {
		opts.GRPCAddress = GRPCAddress
	}
}
