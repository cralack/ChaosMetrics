package config

type Root struct {
	Env       string `mapstructure:"env"`
	DirTree   *DirTree
	Dbconf    *DatabaseConfig `mapstructure:"database"`
	RedisConf *RedisConfig    `mapstructure:"redis"`
	LogConf   *LoggerConfig   `mapstructure:"logger"`
	Fetcher   *FetcherConfig  `mapstructure:"fetcher"`
	System    *SystemConf     `mapstructure:"system"`
}
type DirTree struct {
	WorkDir string
	LogDir  string
	TestDir string
	// ConfDir string
}

func New() *Root {

	return &Root{}
}
