package summoner

import (
	"github.com/gin-gonic/gin"
)

type sumnApi struct{}

func InitSummonerRouter(router *gin.RouterGroup) {
	api := sumnApi{}
	routerPath := ""
	routerGroup := router.Group(routerPath)
	{
		routerGroup.GET("/summoner", api.QuerySummoner)
	}
}
