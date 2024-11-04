package config

type AmqpConfig struct {
	User       string `mapstructure:"user" yaml:"user"`
	Password   string `mapstructure:"password" yaml:"password"`
	Host       string `mapstructure:"host" yaml:"host"`
	Port       string `mapstructure:"port" yaml:"port"`
	URL        string `mapstructure:"url" yaml:"url"`
	Exchange   string `mapstructure:"exchange" yaml:"exchange" default:"amq.direct"`
	Queue      string `mapstructure:"queue" yaml:"queue"`
	RoutingKey string `mapstructure:"routing_key" yaml:"routing_key"`
	AutoDelete bool   `mapstructure:"auto_delete" yaml:"auto_delete"`
}
