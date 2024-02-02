package champion_rank

import (
	"github.com/gin-gonic/gin"
)

type championRankApi struct{}

func InitChampionRankRouter(r *gin.RouterGroup) {
	api := championRankApi{}
	routerPath := ""
	championRouter := r.Group(routerPath)
	{
		championRouter.GET("/ARAM", api.QueryChampionRankARAM)
		championRouter.GET("/CLASSIC", api.QueryChampionRankCLASSIC)
	}
}
