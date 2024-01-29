package champion_detail

import (
	"github.com/gin-gonic/gin"
)

type championDetailApi struct{}

func InitChampionRankRouter(r *gin.Engine) {
	api := championDetailApi{}
	championDetailApi := r.Group("")
	{
		championDetailApi.GET("/champion", api.QueryChampionDetail)
	}
}
