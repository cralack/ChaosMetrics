package champion_rank

import (
	"github.com/cralack/ChaosMetrics/server/app/provider/champion_rank"
	"github.com/cralack/ChaosMetrics/server/model/anres"
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/gin-gonic/gin"
)

type championRankParam struct {
	Loc     string `form:"loc" default:"na1" binding:"required"`        // Region
	Version string `form:"version" default:"14.1.1" binding:"required"` // Version
}

// QueryChampionRankARAM godoc
//
//	@Summary		请求一个ARAM英雄榜
//	@Description	query @version,loc
//	@Accept			application/json
//	@Produce		application/json
//	@Tags			Champion Rank
//	@Param			championRankParam	query		championRankParam	true	"Query champion rank list for aram"
//	@Success		200					{object}	response.Response{data=[]anres.ChampionBrief}
//	@Router			/ARAM [get]
func (a *championRankApi) QueryChampionRankARAM(ctx *gin.Context) {
	var (
		param            championRankParam
		championRankList []*anres.ChampionBrief
		err              error
	)
	championRankService := champion_rank.NewChampionRankService()

	if err = ctx.ShouldBindQuery(&param); err != nil {
		response.FailWithMessage("wrong param", ctx)
		return
	}
	if championRankList, err = championRankService.QueryChampionRank(
		param.Version, param.Loc, "ARAM"); err != nil {
		response.FailWithDetailed(err, "can not find ARAM champion rank list", ctx)
		return
	} else {
		response.OkWithData(championRankList, ctx)
	}
}

// QueryChampionRankCLASSIC godoc
//
//	@Summary		请求一个CLASSIC英雄榜
//	@Description	query @version,loc
//	@Accept			application/json
//	@Produce		application/json
//	@Tags			Champion Rank
//	@Param			championRankParam	query		championRankParam	true	"Query champion rank list for classic"
//	@Success		200					{object}	[]anres.ChampionBrief
//	@Router			/CLASSIC [get]
func (a *championRankApi) QueryChampionRankCLASSIC(ctx *gin.Context) {
	var (
		param            championRankParam
		championRankList []*anres.ChampionBrief
		err              error
	)
	championRankService := champion_rank.NewChampionRankService()

	if err = ctx.ShouldBindQuery(&param); err != nil {
		response.FailWithMessage("wrong param", ctx)
		return
	}
	if championRankList, err = championRankService.QueryChampionRank(
		param.Version, param.Loc, "CLASSIC"); err != nil {
		response.FailWithDetailed(err, "can not find CLASSIC champion rank list", ctx)
		return
	} else {
		response.OkWithData(championRankList, ctx)
	}
}
