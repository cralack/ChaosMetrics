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
//	@Description	query @name,loc
//	@Accept			application/json
//	@Produce		application/json
//	@Tags			Summoner Detail
//	@Param			data	query		summonerQueryParam	true	"Query champion rank list for aram"
//	@Success		200		{object}	response.Response{data=response.SummonerDTO}
//	@Router			/summoner [get]
func (a *sumnApi) QuerySummoner(ctx *gin.Context) {
	var (
		param summonerQueryParam
		sumn  *response.SummonerDTO
		err   error
	)
	sumService := summoner.NewSumnService()
	if err = ctx.ShouldBindQuery(&param); err != nil {
		response.FailWithMessage("wrong param", ctx)
		return
	}
	if sumn = sumService.QuerySummonerByName(param.Name,
		param.Loc); sumn == nil {
		response.FailWithMessage("can not find sumn now,try later", ctx)
		return
	} else {
		response.OkWithData(sumn, ctx)
	}
}
