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
	TEST_ENV Env = iota
	DEV_ENV
	PRODUCT_ENV
)

var (
	GVA_ENV  Env
	GVA_CONF *config.Root
	GVA_DB   *gorm.DB
	GVA_RDB  *redis.Client
	GVA_VP   *viper.Viper
	GVA_LOG  *zap.Logger
)

const (
	WorkerServiceName = "pumper.worker"
	MasterServiceName = "pumper.master"
	TaskPath          = "/tasks"
	ElectPath         = "/resources/election"
)
