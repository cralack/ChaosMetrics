package xviper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cralack/ChaosMetrics/server/global"
	"github.com/cralack/ChaosMetrics/server/internal/config"
	"go.uber.org/zap"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Viper() (*viper.Viper, error) {
	// only init once
	if global.GVA_VP != nil {
		return global.GVA_VP, nil
	}
	v := viper.New()

	// setup dir
	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	conf := global.GVA_CONF
	workDir := curDir[:strings.Index(curDir, "server")+len("server")]
	logDir := filepath.Join(workDir, "log")
	testDir := filepath.Join(workDir, "test")
	if conf.DirTree == nil {
		conf.DirTree = &config.DirTree{
			WorkDir: workDir,
			LogDir:  logDir,
			TestDir: testDir,
		}
	}
	v.AddConfigPath(workDir)

	// setup config file name
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// read conf file
	err = v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// Watching and re-reading config files
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		// handler func
		global.GVA_LOG.Info("config file changed:",
			zap.String("filename", e.Name))
		if err = v.Unmarshal(conf); err != nil {
			panic(err)
		}
	})

	if err = v.Unmarshal(conf); err != nil {
		panic(err)
	}
	switch conf.Env {
	case "test":
		global.GVA_ENV = global.TEST_ENV
	case "dev":
		global.GVA_ENV = global.DEV_ENV
	case "product":
		global.GVA_ENV = global.PRODUCT_ENV
	}
	return v, nil
}
