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
	ChaEnv    Env
	ChaConf   *config.Root
	ChaDB     *gorm.DB
	ChaRDB    *redis.Client
	ChaViper  *viper.Viper
	ChaLogger *zap.Logger
)

const (
	WorkerServiceName = "pumper.worker"
	MasterServiceName = "pumper.master"
	TaskPath          = "/tasks"
	ElectPath         = "/resources/election"
)
