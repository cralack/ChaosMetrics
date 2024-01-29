package champion_detail

import (
	"sort"

	"github.com/cralack/ChaosMetrics/server/app/provider/champion_detail"
	"github.com/cralack/ChaosMetrics/server/model/anres"
	"github.com/gin-gonic/gin"
)

type championDetailParam struct {
	Name    string `form:"name" default:"Ahri" binding:"required"`      // Champion name
	Loc     string `form:"loc" default:"na1" binding:"required"`        // Region
	Version string `form:"version" default:"14.1.1" binding:"required"` // Version
	Mode    string `form:"mode" default:"CLASSIC" binding:"required"`   // Game mode
}

// QueryChampionDetail godoc
//
//	@Summary		请求一个英雄详情
//	@Description	请求一个英雄详情 @name,version,loc,mode
//	@Accept			application/json
//	@Produce		application/json
//	@Tags			Champion Detail
//	@Param			championDetailParam	query		championDetailParam	true	"Query champion rank list for aram"
//	@Success		200					{object}	anres.ChampionDetail
//	@Router			/champion [get]
func (c *championDetailApi) QueryChampionDetail(ctx *gin.Context) {
	var (
		param    championDetailParam
		champion *anres.ChampionDetail
		err      error
	)
	championRankService := champion_detail.NewChampionDetailService()
	if err = ctx.ShouldBindQuery(&param); err != nil {
		ctx.JSON(400, gin.H{
			"msg": "wrong param",
		})
		return
	}
	if champion, err = championRankService.QueryChampionDetail(
		param.Name, param.Version, param.Loc, param.Mode); err != nil {
		ctx.JSON(400, gin.H{
			"msg": "can not find champion detail by " + err.Error(),
		})
		return
	} else {
		for key := range champion.ItemWin {
			champion.ItemWin[key] = ShrinkMap(champion.ItemWin[key])
		}
		champion.PerkWin = ShrinkMap(champion.PerkWin)

		ctx.JSON(200, gin.H{
			"msg":  "query success",
			"data": champion,
		})
	}
}

func ShrinkMap(src map[string]int) (des map[string]int) {
	if len(src) < 10 {
		return
	}
	type node struct {
		data string
		idx  int
	}
	tmp := make([]*node, 0, len(src))
	for k, v := range src {
		tmp = append(tmp, &node{
			data: k,
			idx:  v,
		})
	}
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].idx > tmp[j].idx
	})
	tmp = tmp[:10]
	des = make(map[string]int)
	for _, n := range tmp {
		des[n.data] = n.idx
	}

	return
}
