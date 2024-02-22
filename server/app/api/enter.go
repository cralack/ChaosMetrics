package api

import (
	"net/http"

	"github.com/cralack/ChaosMetrics/server/app/api/champion_detail"
	"github.com/cralack/ChaosMetrics/server/app/api/champion_rank"
	"github.com/cralack/ChaosMetrics/server/app/api/comment"
	"github.com/cralack/ChaosMetrics/server/app/api/common"
	"github.com/cralack/ChaosMetrics/server/app/api/item"
	"github.com/cralack/ChaosMetrics/server/app/api/summoner"
	"github.com/cralack/ChaosMetrics/server/app/api/user"
	"github.com/cralack/ChaosMetrics/server/app/middleware"
	"github.com/cralack/ChaosMetrics/server/docs"
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine) {
	// register swagger
	docs.SwaggerInfo.BasePath = global.ChaConf.System.RouterPrefix
	r.GET(global.ChaConf.System.RouterPrefix+"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler))
	global.ChaLogger.Info("register swagger handler")
	// cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, global.TokenKey)
	corsConfig.AllowAllOrigins = true
	c := cors.New(corsConfig)
	r.Use(c)

	PublicGroup := r.Group(global.ChaConf.System.RouterPrefix)
	{
		PublicGroup.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, "ok")
		})
		item.InitItemRouter(PublicGroup)
		champion_rank.InitChampionRankRouter(PublicGroup)
		champion_detail.InitChampionRankRouter(PublicGroup)
		summoner.InitSummonerRouter(PublicGroup)
		user.InitUserRouter(PublicGroup)
		common.InitCommonAPI(PublicGroup)
	}

	PrivateGroup := r.Group(global.ChaConf.System.RouterPrefix)
	PrivateGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	{
		comment.InitCommentRouter(PrivateGroup)
	}

	global.ChaLogger.Info("router register success")
}
