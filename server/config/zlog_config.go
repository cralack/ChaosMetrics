package config

type LoggerConfig struct {
	MaxSize    int  `mapstructure:"maxsize" yaml:"maxsize"`       // 日志的最大大小，以M为单位
	MaxBackups int  `mapstructure:"maxbackups" yaml:"maxbackups"` // 保留的旧日志文件的最大数量
	MaxAge     int  `mapstructure:"maxage" yaml:"maxage"`         // 保留旧日志文件的最大天数
	LocalTime  bool `mapstructure:"localtime" yaml:"localtime"`   // 是否使用本地时间
	Compress   bool `mapstructure:"compress" yaml:"compress"`     // 是否压缩旧日志文件
}

func NewLoggerConfig() *LoggerConfig {
	return &LoggerConfig{}
}
