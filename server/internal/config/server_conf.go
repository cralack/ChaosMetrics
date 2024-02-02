package config

import (
	"time"
)

type SystemConf struct {
	ID                string        `mapstructure:"id" yaml:"id"`                                   // 服务ID
	Name              string        `mapstructure:"name" yaml:"name"`                               // 服务名
	GRPCListenAddress string        `mapstructure:"grpc_listen_address" yaml:"grpc_listen_address"` // gRPC监听地址
	HTTPListenAddress string        `mapstructure:"http_listen_address" yaml:"http_listen_address"` // HTTP监听地址
	RegistryAddress   string        `mapstructure:"registry_address" yaml:"registry_address"`       // 注册中心地址
	RegisterTTL       time.Duration `mapstructure:"register_ttl" yaml:"register_ttl"`               // 注册TTL（Time-to-Live）
	RegisterInterval  time.Duration `mapstructure:"register_interval" yaml:"register_interval"`     // 注册间隔
	ClientTimeOut     time.Duration `mapstructure:"client_timeout" yaml:"client_timeout"`           // 客户端超时时间
	RouterPrefix      string        `mapstructure:"router_prefix" yaml:"router_prefix"`             // 路由全局前缀
}
