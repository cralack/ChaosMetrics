package config

import (
	"time"
)

type ServerConfig struct {
	ID               string        `mapstructure:"id" yaml:"id"`                               // 服务ID
	RegistryAddress  string        `mapstructure:"registry_address" yaml:"registry_address"`   // 注册中心地址
	RegisterTTL      time.Duration `mapstructure:"register_ttl" yaml:"register_ttl"`           // 注册TTL（Time-to-Live）
	RegisterInterval time.Duration `mapstructure:"register_interval" yaml:"register_interval"` // 注册间隔
	ClientTimeOut    time.Duration `mapstructure:"client_timeout" yaml:"client_timeout"`       // 客户端超时时间
}
