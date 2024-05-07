package hero_data

import (
	"sort"

	"github.com/cralack/ChaosMetrics/server/app/provider/hero_data"
	"github.com/cralack/ChaosMetrics/server/model/anres"
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/gin-gonic/gin"
)

type heroDataParam struct {
	Name    string `form:"name" default:"Ahri" binding:"required"`      // Champion name
	Loc     string `form:"loc" default:"na1" binding:"required"`        // Region
	Version string `form:"version" default:"14.1.1" binding:"required"` // Version
	Mode    string `form:"mode" default:"CLASSIC" binding:"required"`   // Game mode
}

// QueryHeroData godoc
//
//	@Summary		请求一个英雄详情
//	@Description	query @name,version,loc,mode
//	@Accept			application/json
//	@Produce		application/json
//	@Tags			Hero Data
//	@Param			key	query		heroDataParam	true	"Query champion rank list for aram"
//	@Success		200		{object}	response.Response{key=anres.ChampionDetail}
//	@Router			/hero [get]
func (a *heroDataApi) QueryHeroData(ctx *gin.Context) {
	var (
		param    heroDataParam
		champion *anres.ChampionDetail
		err      error
	)
	heroDataServ := hero_data.NewHeroDataService()
	if err = ctx.ShouldBindQuery(&param); err != nil {
		response.FailWithMessage("wrong param", ctx)
		return
	}
	if champion, err = heroDataServ.QueryHeroData(
		param.Name, param.Version, param.Loc, param.Mode); err != nil {
		response.FailWithDetailed(err, "can not find champion", ctx)
		return
	} else {
		for key := range champion.ItemWin {
			champion.ItemWin[key] = ShrinkMap(champion.ItemWin[key])
		}
		champion.PerkWin = ShrinkMap(champion.PerkWin)
		champion.SkillWin = ShrinkMap(champion.SkillWin)
		response.OkWithQuiet(champion, ctx)
	}
}

func ShrinkMap(src map[string]*anres.Stats) (des map[string]*anres.Stats) {
	if len(src) < 10 {
		return src
	}
	type node struct {
		key  string
		stat *anres.Stats
	}
	tmp := make([]*node, 0, len(src))
	for k, v := range src {
		tmp = append(tmp, &node{
			key:  k,
			stat: v,
		})
	}
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].stat.Wins > tmp[j].stat.Wins
	})
	tmp = tmp[:10]
	des = map[string]*anres.Stats{}
	for _, n := range tmp {
		des[n.key] = n.stat
	}
	return
}
