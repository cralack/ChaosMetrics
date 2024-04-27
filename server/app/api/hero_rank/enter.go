package hero_rank

import (
	"github.com/gin-gonic/gin"
)

type heroRankApi struct{}

func InitHeroRankRouter(r *gin.RouterGroup) {
	api := heroRankApi{}
	routerPath := ""
	routerGroup := r.Group(routerPath)
	{
		routerGroup.GET("/ARAM", api.QueryHeroRankARAM)
		routerGroup.GET("/CLASSIC", api.QueryHeroRankCLASSIC)
	}
}
