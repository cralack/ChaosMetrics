package config

type Server struct {
	Env     string `mapstructure:"env"`
	DirTree DirTree
	Dbconf  *DatabaseConfig `mapstructure:"database"`
	LogConf *LoggerConfig   `mapstructure:"logger"`
	Fetcher *FetcherConfig  `mapstructure:"fetcher"`
}
type DirTree struct {
	WorkDir string
	LogDir  string
	TestDir string

	// ConfDir string
}

func New() *Server {

	return &Server{}
}
