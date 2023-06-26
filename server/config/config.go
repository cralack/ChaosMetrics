package config

type Server struct {
	DirTree DirTree
	Dbconf  *DatabaseConfig
	LogConf *LoggerConfig
}
type DirTree struct {
	WordDir string
	ConfDir string
	LogDIr  string
}

func New() *Server {
	return &Server{
		Dbconf:  NewDBConfig(),
		LogConf: NewLoggerConfig(),
	}
}
