package champion_rank

import (
	"github.com/cralack/ChaosMetrics/server/app/provider/champion_rank"
	"github.com/cralack/ChaosMetrics/server/model/anres"
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
//	@Success		200					{object}	[]anres.ChampionBrief
//	@Router			/ARAM [get]
func (c *championRankApi) QueryChampionRankARAM(ctx *gin.Context) {
	var (
		param            championRankParam
		championRankList []*anres.ChampionBrief
		err              error
	)
	championRankService := champion_rank.NewChampionRankService()

	if err = ctx.ShouldBindQuery(&param); err != nil {
		ctx.JSON(400, gin.H{
			"msg": "wrong param",
		})
		return
	}
	if championRankList, err = championRankService.QueryChampionRank(
		param.Version, param.Loc, "ARAM"); err != nil {
		ctx.JSON(400, gin.H{
			"msg": "can not find champion rank list by " + err.Error(),
		})
		return
	} else {
		ctx.JSON(200, gin.H{
			"msg":  "query success",
			"data": championRankList,
		})
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
func (c *championRankApi) QueryChampionRankCLASSIC(ctx *gin.Context) {
	var (
		param            championRankParam
		championRankList []*anres.ChampionBrief
		err              error
	)
	championRankService := champion_rank.NewChampionRankService()

	if err = ctx.ShouldBindQuery(&param); err != nil {
		ctx.JSON(400, gin.H{
			"msg": "wrong param",
		})
		return
	}
	if championRankList, err = championRankService.QueryChampionRank(
		param.Version, param.Loc, "CLASSIC"); err != nil {
		ctx.JSON(400, gin.H{
			"msg": "can not find champion rank list by " + err.Error(),
		})
		return
	} else {
		ctx.JSON(200, gin.H{
			"msg":  "query success",
			"data": championRankList,
		})
	}
}
