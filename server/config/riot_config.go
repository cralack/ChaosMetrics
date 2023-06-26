package config

type RiotConf struct {
	Apikey      string `mapstructure:"api_key" yaml:"api_key"` // Riot API key
	RateLimiter struct {
		EachSec  int `mapstructure:"each_sec" yaml:"each_sec"`   // Requests per second
		Each2min int `mapstructure:"each_2min" yaml:"each_2min"` // Requests per 2 minutes
	} `mapstructure:"rate_limiter" yaml:"rate_limiter"` // Rate limiter settings
}
