package config

type Root struct {
	Env       string `mapstructure:"env"`
	DirTree   *DirTree
	Dbconf    *DatabaseConfig `mapstructure:"database"`
	RedisConf *RedisConfig    `mapstructure:"redis"`
	LogConf   *LoggerConfig   `mapstructure:"logger"`
	Fetcher   *FetcherConfig  `mapstructure:"fetcher"`
	Micro     *MicroServ      `mapstructure:"micro"`
	Router    *Router         `mapstructure:"router"`
	JwtConf   *JwtConfig      `mapstructure:"jwt"`
	EmailConf *EmailConf      `mapstructure:"smtp"`
	AmqpConf  *AmqpConfig     `mapstructure:"amqp"`
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
