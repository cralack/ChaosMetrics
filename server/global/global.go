package global

import (
	"github.com/cralack/ChaosMetrics/server/config"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	TEST_ENV = iota
	DEV_ENV
	PRODUCT_ENV
)

var (
	GVA_ENV  uint
	GVA_CONF *config.Server
	GVA_DB   *gorm.DB
	GVA_VP   *viper.Viper
	GVA_LOG  *zap.Logger
)
