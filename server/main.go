package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

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
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// timer := time.NewTimer(time.Second * 5)
	errChan := make(chan error, 1)
	go func() {
		if err := cmd.RunCmd(ctx); err != nil {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		cancel()
		global.ChaLogger.Error("run app failed:", zap.Error(err))
	case sig := <-quit:
		cancel()
		global.ChaLogger.Info("received signal:", zap.String("signal", sig.String()))
	}
}
