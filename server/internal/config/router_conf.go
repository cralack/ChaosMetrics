package config

type Router struct {
	Domain       string `mapstructure:"domain" yaml:"domain"`
	DomainHost   string `mapstructure:"domain_host" yaml:"domain_host"`
	DomainPort   string `mapstructure:"domain_port" yaml:"domain_port"`
	RouterPrefix string `mapstructure:"router_prefix" yaml:"router_prefix"` // 路由全局前缀
}
