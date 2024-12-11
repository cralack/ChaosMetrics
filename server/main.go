package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cralack/ChaosMetrics/server/cmd"
	_ "github.com/cralack/ChaosMetrics/server/init"
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"go.uber.org/zap"
)

//	@title						ChaosMetrics API接口文档
//	@version					0.9
//	@description				使用Riot官方API获取数据进行分析、统计项目
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//	@host						localhost:8080
//	@BasePath					/
//	@securityDefinitions.apikey	TokenAuth
//	@in							header
//	@name						x-token
//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
		time.Sleep(time.Second * 1)
	}()

	go func() {
		if err := cmd.RunCmd(ctx); err != nil {
			global.ChaLogger.Fatal("run app failed:", zap.Error(err))
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit // 阻塞等待接收 channel 数据
	global.ChaLogger.Info("exiting...")
}
