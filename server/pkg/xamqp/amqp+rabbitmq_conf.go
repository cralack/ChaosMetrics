package xamqp

import (
	"context"
)

type Config struct {
	Context    context.Context
	Addr       string
	Exchange   string
	Queue      string
	RoutingKey string
	AutoDelete bool
}

type Setup func(*Config)

var defaultConfig = Config{
	Addr:       "",
	Exchange:   Exchange,
	Queue:      Queue,
	RoutingKey: RoutingKey,
	AutoDelete: false,
}

func WithAddr(addr string) Setup {
	return func(opts *Config) {
		opts.Addr = addr
	}
}

func WithExchange(exchange string) Setup {
	return func(opts *Config) {
		opts.Exchange = exchange
	}
}

func WithQueue(queue string) Setup {
	return func(opts *Config) {
		opts.Queue = queue
	}
}

func WithRoutingKey(routingKey string) Setup {
	return func(opts *Config) {
		opts.RoutingKey = routingKey
	}
}

func WithAutodelete(autodelete bool) Setup {
	return func(opts *Config) {
		opts.AutoDelete = autodelete
	}
}

func WithContext(ctx context.Context) Setup {
	return func(opts *Config) {
		opts.Context = ctx
	}
}
