package config

type DatabaseConfig struct {
	DSN          string `mapstructure:"dsn" yaml:"dsn"`                     // 数据源名称
	Driver       string `mapstructure:"driver" yaml:"driver"`               // 数据库驱动名称
	Host         string `mapstructure:"host" yaml:"host"`                   // 数据库主机
	Port         int    `mapstructure:"port" yaml:"port"`                   // 数据库端口
	DBName       string `mapstructure:"dbname" yaml:"dbname"`               // 数据库名称
	Username     string `mapstructure:"username" yaml:"username"`           // 数据库用户名
	Password     string `mapstructure:"password" yaml:"password"`           // 数据库密码
	Charset      string `mapstructure:"charset" yaml:"charset"`             // 数据库字符集
	Collation    string `mapstructure:"collation" yaml:"collation"`         // 数据库排序规则
	Timeout      string `mapstructure:"timeout" yaml:"timeout"`             // 连接超时时间
	ReadTimeout  string `mapstructure:"read_timeout" yaml:"read_timeout"`   // 读取超时时间
	WriteTimeout string `mapstructure:"write_timeout" yaml:"write_timeout"` // 写入超时时间
	ParseTime    bool   `mapstructure:"parse_time" yaml:"parse_time"`       // 解析时间戳
	Protocol     string `mapstructure:"protocol" yaml:"protocol"`           // 数据库协议
	Loc          string `mapstructure:"loc" yaml:"loc"`                     // 时区

	ConnMaxIdle     int    `mapstructure:"conn_max_idle" yaml:"conn_max_idle"`         // 连接池最大空闲连接数
	ConnMaxOpen     int    `mapstructure:"conn_max_open" yaml:"conn_max_open"`         // 连接池最大打开连接数
	ConnMaxLifetime string `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime"` // 连接池连接的最大生命周期
}

func NewDBConfig() *DatabaseConfig {
	return &DatabaseConfig{}
}
