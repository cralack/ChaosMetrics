package hero_detail

import (
	"github.com/cralack/ChaosMetrics/server/app/provider/hero_detail"
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/gin-gonic/gin"
)

type heroDetailParam struct {
	Name    string `form:"name" default:"Ahri" binding:"required"`      // Champion name
	Version string `form:"version" default:"14.1.1" binding:"required"` // Version
	Lang    string `form:"lang"  default:"zh_CN" binding:"required"`
}

// QueryHeroDetail godoc
//
//	@Summary		请求一个英雄详细资料
//	@Description	query @name,version,lang
//	@Accept			application/json
//	@Produce		application/json
//	@Tags			Hero Detail
//	@Param			data	query		heroDetailParam	true	"Query hero detail"
//	@Success		200		{object}	response.Response{data}
//	@Router			/hero_detail [get]
func (a *heroDetailApi) QueryHeroDetail(ctx *gin.Context) {
	var (
		param heroDetailParam
		hero  *riotmodel.ChampionDTO
		err   error
	)
	heroDetailServ := hero_detail.NewHeroDetailService()

	if err = ctx.ShouldBindQuery(&param); err != nil {
		response.FailWithMessage("wrong param", ctx)
		return
	}
	if hero, err = heroDetailServ.QueryHeroDetail(
		param.Name, param.Version, param.Lang); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	} else {
		response.OkWithQuiet(hero, ctx)
	}
}
