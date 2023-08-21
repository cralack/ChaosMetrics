package config

import (
	"time"
)

type RedisConfig struct {
	Host            string        `yaml:"host" mapstructure:"host"`                           // Redis 服务器地址
	Port            int           `yaml:"port" mapstructure:"port"`                           // Redis 端口
	DB              int           `yaml:"db" mapstructure:"db"`                               // Redis 数据库索引
	Username        string        `yaml:"username" mapstructure:"username"`                   // Redis 用户名
	Password        string        `yaml:"password" mapstructure:"password"`                   // Redis 密码
	Timeout         time.Duration `yaml:"timeout" mapstructure:"timeout"`                     // 连接超时时间
	ReadTimeout     time.Duration `yaml:"read_timeout" mapstructure:"read_timeout"`           // 读超时时间
	WriteTimeout    time.Duration `yaml:"write_timeout" mapstructure:"write_timeout"`         // 写超时时间
	ConnMinIdle     int           `yaml:"conn_min_idle" mapstructure:"conn_min_idle"`         // 连接池最小空闲连接数
	ConnMaxOpen     int           `yaml:"conn_max_open" mapstructure:"conn_max_open"`         // 连接池最大连接数
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"` // 连接数最大生命周期
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idletime" mapstructure:"conn_max_idletime"` // 连接数空闲时长
}
