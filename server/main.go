package main

import (
	"github.com/cralack/ChaosMetrics/server/cmd"
	_ "github.com/cralack/ChaosMetrics/server/init"
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
	if err := cmd.RunCmd(); err != nil {
		panic(err)
	}
}
