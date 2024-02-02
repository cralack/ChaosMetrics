package champion_detail

import (
	"github.com/gin-gonic/gin"
)

type championDetailApi struct{}

func InitChampionRankRouter(r *gin.RouterGroup) {
	api := championDetailApi{}
	routerPath := ""
	championDetailApi := r.Group(routerPath)
	{
		championDetailApi.GET("/champion", api.QueryChampionDetail)
	}
}
