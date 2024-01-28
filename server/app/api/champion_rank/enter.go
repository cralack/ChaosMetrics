package champion_rank

import (
	"github.com/gin-gonic/gin"
)

type championRankApi struct{}

func InitChampionRankRouter(r *gin.Engine) {
	api := championRankApi{}
	championRankApi := r.Group("/champion")
	{
		championRankApi.GET("/ARAM", api.QueryChampionRankARAM)
		championRankApi.GET("/CLASSIC", api.QueryChampionRankCLASSIC)
	}
}
