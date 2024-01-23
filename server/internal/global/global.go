package global

import (
	"github.com/cralack/ChaosMetrics/server/internal/config"
	"github.com/redis/go-redis/v9"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Env int

const (
	TestEnv Env = iota
	DevEnv
	ProductEnv
)

var (
	GvaEnv  Env
	GvaConf *config.Root
	GvaDb   *gorm.DB
	GvaRdb  *redis.Client
	GvaVp   *viper.Viper
	GvaLog  *zap.Logger
)

const (
	WorkerServiceName = "pumper.worker"
	MasterServiceName = "pumper.master"
	TaskPath          = "/tasks"
	ElectPath         = "/resources/election"
)
