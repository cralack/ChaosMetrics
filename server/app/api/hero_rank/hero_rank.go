package hero_rank

import (
	"github.com/cralack/ChaosMetrics/server/app/provider/hero_rank"
	"github.com/cralack/ChaosMetrics/server/model/anres"
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/gin-gonic/gin"
)

type heroRankParam struct {
	Loc     string `form:"loc" default:"na1" binding:"required"`        // Region
	Version string `form:"version" default:"14.1.1" binding:"required"` // Version
}

// QueryHeroRankARAM godoc
//
//	@Summary		请求一个ARAM英雄榜
//	@Description	query @version,loc
//	@Accept			application/json
//	@Produce		application/json
//	@Tags			Hero Rank
//	@Param			data	query		heroRankParam	true	"Query champion rank list for aram"
//	@Success		200		{object}	response.Response{data=[]anres.ChampionBrief}
//	@Router			/ARAM [get]
func (a *heroRankApi) QueryHeroRankARAM(ctx *gin.Context) {
	var (
		param            heroRankParam
		championRankList []*anres.ChampionBrief
		err              error
	)
	heroRankServ := hero_rank.NewHeroRankService()

	if err = ctx.ShouldBindQuery(&param); err != nil {
		response.FailWithMessage("wrong param", ctx)
		return
	}
	if championRankList, err = heroRankServ.QueryHeroRank(
		param.Version, param.Loc, "ARAM"); err != nil {
		response.FailWithDetailed(err, "can not find ARAM champion rank list", ctx)
		return
	} else {
		response.OkWithQuiet(championRankList, ctx)
	}
}

// QueryHeroRankCLASSIC godoc
//
//	@Summary		请求一个CLASSIC英雄榜
//	@Description	query @version,loc
//	@Accept			application/json
//	@Produce		application/json
//	@Tags			Hero Rank
//	@Param			heroRankParam	query		heroRankParam	true	"Query champion rank list for classic"
//	@Success		200				{object}	[]anres.ChampionBrief
//	@Router			/CLASSIC [get]
func (a *heroRankApi) QueryHeroRankCLASSIC(ctx *gin.Context) {
	var (
		param            heroRankParam
		championRankList []*anres.ChampionBrief
		err              error
	)
	heroRankServ := hero_rank.NewHeroRankService()

	if err = ctx.ShouldBindQuery(&param); err != nil {
		response.FailWithMessage("wrong param", ctx)
		return
	}
	if championRankList, err = heroRankServ.QueryHeroRank(
		param.Version, param.Loc, "CLASSIC"); err != nil {
		response.FailWithDetailed(err, "can not find CLASSIC champion rank list", ctx)
		return
	} else {
		response.OkWithQuiet(championRankList, ctx)
	}
}
