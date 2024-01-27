package route

import (
	"net/http"

	"github.com/cralack/ChaosMetrics/server/app/api/item"
	"github.com/cralack/ChaosMetrics/server/docs"
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegiserRoutes(r *gin.Engine) {
	docs.SwaggerInfo.BasePath = global.ChaConf.System.RouterPrefix
	r.GET(global.ChaConf.System.RouterPrefix+"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler))
	global.ChaLogger.Info("register swagger handler")

	PublicGroup := r.Group(global.ChaConf.System.RouterPrefix)
	{
		PublicGroup.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, "ok")
		})
		item.InitItemRouter(r)
	}
	global.ChaLogger.Info("router register success")
}
