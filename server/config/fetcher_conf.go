package config

import "time"

type FetcherConfig struct {
	Timeout time.Duration `mapstructure:"timeout" yaml:"timeout"` // 连接超时时间

	// Fetcher的请求头配置
	HeaderConfig struct {
		UserAgent      string `mapstructure:"user_agent" yaml:"user_agent"`           // 用户代理
		AcceptLanguage string `mapstructure:"accept_language" yaml:"accept_language"` // 接受的语言
		AcceptCharset  string `mapstructure:"accept_charset" yaml:"accept_charset"`   // 接受的字符集
		Origin         string `mapstructure:"origin" yaml:"origin"`                   // 请求来源
		XRiotToken     string `mapstructure:"x_riot_token" yaml:"x_riot_token"`       // Riot API的密钥
	} `mapstructure:"header" yaml:"header"`

	// Fetcher的速率限制配置
	RateLimiterConfig struct {
		EachSec  int `mapstructure:"each_sec" yaml:"each_sec"`   // 每秒允许的API请求次数限制
		Each2Min int `mapstructure:"each_2min" yaml:"each_2min"` // 每两分钟允许的API请求次数限制
	} `mapstructure:"rate_limiter" yaml:"rate_limiter"`
}
