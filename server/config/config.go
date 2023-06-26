package config

type Server struct {
	DirTree DirTree
	Dbconf  *DatabaseConfig `mapstructure:"database"`
	LogConf *LoggerConfig   `mapstructure:"logger"`
	Riot    *RiotConf       `mapstructure:"riot"`
}
type DirTree struct {
	WordDir string
	ConfDir string
	LogDIr  string
}

func New() *Server {
	return &Server{
		Dbconf:  &DatabaseConfig{},
		LogConf: &LoggerConfig{},
	}
}
