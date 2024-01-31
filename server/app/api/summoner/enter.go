package summoner

import (
	"github.com/gin-gonic/gin"
)

type sumnApi struct{}

func InitSummonerRouter(r *gin.Engine) {
	api := sumnApi{}
	SummonerApi := r.Group("")
	{
		SummonerApi.GET("/summoner", api.QuerySummoner)
	}
}
