package vconf

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"ChaosMetrics/server/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Viper() (*viper.Viper, error) {
	//only init once
	if global.GVA_VP != nil {
		return global.GVA_VP, nil
	}
	v := viper.New()

	//setup config file name
	//todo:diff conf file
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	//setup dir
	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	workDir := curDir[:strings.Index(curDir, "server")+len("server")]
	confDir := filepath.Join(workDir, "config")
	logDir := filepath.Join(workDir, "log")
	if global.GVA_CONF.DirTree.WordDir == "" {
		global.GVA_CONF.DirTree.WordDir = workDir
		global.GVA_CONF.DirTree.LogDIr = logDir
	}
	v.AddConfigPath(confDir)

	//
	err = v.ReadInConfig()
	if err != nil {
		log.Fatal("Error reading config file")
	}

	//Watching and re-reading config files
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		//handler func
		// global.GVA_LOG.Info("Config file changed:",
		// 	zap.String("conf", e.Name))
	})

	return v, nil
}
