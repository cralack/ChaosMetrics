package config

type EmailConf struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     int    `mapstructure:"port" yaml:"port"`
	From     string `mapstructure:"from" yaml:"from"`
	Username string `mapstructure:"username" yaml:"username"`
	Passwd   string `mapstructure:"password" yaml:"password"`
}
