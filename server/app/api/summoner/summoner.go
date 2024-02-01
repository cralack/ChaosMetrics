package summoner

import (
	"github.com/cralack/ChaosMetrics/server/app/provider/summoner"
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/gin-gonic/gin"
)

type summonerQueryParam struct {
	Name string `form:"name" default:"pwVx hysamirapwd" binding:"required"` // Summoner name
	Loc  string `form:"loc" default:"na1" binding:"required"`               // Region
}

// QuerySummoner godoc
//
//	@Summary		请求一个召唤师详情
//	@Description	请求一个召唤师详情 @name,loc
//	@Accept			application/json
//	@Produce		application/json
//	@Tags			Summoner Detail
//	@Param			summonerQueryParam	query		summonerQueryParam	true	"Query champion rank list for aram"
//	@Success		200					{object}	response.SummonerDTO
//	@Router			/summoner [get]
func (a *sumnApi) QuerySummoner(ctx *gin.Context) {
	var (
		param summonerQueryParam
		sumn  *response.SummonerDTO
		err   error
	)
	sumService := summoner.NewSumnService()
	if err = ctx.ShouldBindQuery(&param); err != nil {
		ctx.JSON(400, gin.H{
			"msg": "wrong param",
		})
		return
	}
	if sumn = sumService.QuerySummonerByName(param.Name,
		param.Loc); sumn == nil {
		ctx.JSON(404, gin.H{
			"msg": "can not find sumn now,try later",
		})
		return
	} else {
		ctx.JSON(200, gin.H{
			"msg":  "query success",
			"data": sumn,
		})
	}
}
