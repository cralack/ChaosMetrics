package config

type ServerConfig struct {
	GRPCListenAddress string `mapstructure:"grpc_listen_address" yaml:"grpc_listen_address"` // gRPC监听地址
	HTTPListenAddress string `mapstructure:"http_listen_address" yaml:"http_listen_address"` // HTTP监听地址
	ID                string `mapstructure:"id" yaml:"id"`                                   // 服务ID
	RegistryAddress   string `mapstructure:"registry_address" yaml:"registry_address"`       // 注册中心地址
	RegisterTTL       int    `mapstructure:"register_ttl" yaml:"register_ttl"`               // 注册TTL（Time-to-Live）
	RegisterInterval  int    `mapstructure:"register_interval" yaml:"register_interval"`     // 注册间隔
	Name              string `mapstructure:"name" yaml:"name"`                               // 服务名称
	ClientTimeOut     int    `mapstructure:"client_timeout" yaml:"client_timeout"`           // 客户端超时时间
}
