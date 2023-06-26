package vconf

import (
	"fmt"
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

	//setup dir
	curDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	workDir := curDir[:strings.Index(curDir, "server")+len("server")]
	logDir := filepath.Join(workDir, "log")
	if global.GVA_CONF.DirTree.WordDir == "" {
		global.GVA_CONF.DirTree.WordDir = workDir
		global.GVA_CONF.DirTree.LogDIr = logDir
	}
	v.AddConfigPath(workDir)

	//setup config file name
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	//read conf file
	err = v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	//Watching and re-reading config files
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		//handler func
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(global.GVA_CONF); err != nil {
			panic(err)
		}
	})

	if err = v.Unmarshal(global.GVA_CONF); err != nil {
		panic(err)
	}
	switch global.GVA_CONF.Env {
	case "test":
		global.GVA_ENV = global.TEST_ENV
	case "dev":
		global.GVA_ENV = global.DEV_ENV
	case "product":
		global.GVA_ENV = global.PRODUCT_ENV
	}
	return v, nil
}
