package summoner

import (
	"github.com/gin-gonic/gin"
)

type sumnApi struct{}

func InitSummonerRouter(router *gin.RouterGroup) {
	api := sumnApi{}
	routerPath := ""
	summonerRouter := router.Group(routerPath)
	{
		summonerRouter.GET("/summoner", api.QuerySummoner)
	}
}
