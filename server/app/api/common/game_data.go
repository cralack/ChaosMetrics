package common

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/gin-gonic/gin"
)

// GetGameVersions 获取当前大版本号下的所有版本号。
//
//	@Summary		请求当前大版本号下的版本列表
//	@Description	获取与当前大版本号相关的版本列表
//	@Accept			application/json
//	@Produce		application/json
//	@Tags			Common Game Data
//	@Success		200		{object}	response.Response{data=[]string}
//	@Failure		default	{object}	response.Response
//	@Router			/gameversion [get]
func (a *cmnApi) GetGameVersions(ctx *gin.Context) {
	res := utils.GetCurMajorVersions()
	if len(res) == 0 {
		response.Fail(ctx)
		return
	}

	response.OkWithQuiet(res, ctx)
}

type QueryParam struct {
	Version string `form:"version" default:"14.5.1"` // Version
	Lang    string `form:"lang"  default:"zh_CN" binding:"required"`
}

// GetPerksData 依据游戏版本和语言获取符文数据。
//
//	@Summary		请求特定版本和语言的符文数据
//	@Description	根据提供的版本和语言信息，查询并返回符文数据
//	@Accept			application/json
//	@Produce		application/json
//	@Tags			Common Game Data
//	@Param			data	query		QueryParam	true	"Query Perks"
//	@Success		200		{object}	response.Response{data=[]riotmodel.Perk}
//	@Failure		default	{object}	response.Response
//	@Router			/perks [get]
func (a *cmnApi) GetPerksData(ctx *gin.Context) {
	perks := make([]*riotmodel.Perk, 0, 5)
	keys := []string{"8000", "8100", "8200", "8300", "8400"}
	var param QueryParam
	if err := ctx.ShouldBindQuery(&param); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	vidx, _ := utils.ConvertVersionToIdx(param.Version)
	for i := range keys {
		keys[i] = fmt.Sprintf("%s@%d", keys[i], vidx)
	}
	values := global.ChaRDB.HMGet(context.Background(), "/perks/"+param.Lang, keys...).Val()
	for _, v := range values {
		if v == nil {
			continue
		}
		var perk *riotmodel.Perk
		if err := json.Unmarshal([]byte(v.(string)), &perk); err != nil {
			continue
		}
		perks = append(perks, perk)
	}
	if len(perks) != 5 {
		response.FailWithMessage("get perks failed", ctx)
		return
	}

	response.OkWithQuiet(perks, ctx)
}

// GetSpellData 依据游戏版本和语言获取召唤师技能数据。
//
//	@Summary		请求特定版本和语言的召唤师技能
//	@Description	根据提供的版本和语言信息，查询并返回召唤师技能
//	@Accept			application/json
//	@Produce		application/json
//	@Tags			Common Game Data
//	@Param			data	query		QueryParam	true	"Query Spells"
//	@Success		200		{object}	response.Response{data}
//	@Failure		default	{object}	response.Response
//	@Router			/spells [get]
func (a *cmnApi) GetSpellData(ctx *gin.Context) {
	var (
		param  QueryParam
		res    *riotmodel.SpellList
		spells = make([]*riotmodel.SummonerSpell, 0)
	)

	if err := ctx.ShouldBindQuery(&param); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	vidx, _ := utils.ConvertVersionToIdx(param.Version)
	values := global.ChaRDB.HGet(context.Background(), "/spells/", fmt.Sprintf("%d@%s", vidx, param.Lang)).Val()
	if err := json.Unmarshal([]byte(values), &res); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	for _, v := range res.Data {
		v.Effect = nil
		v.EffectBurn = nil
		v.Vars = nil
		spells = append(spells, v)
	}
	response.OkWithQuiet(spells, ctx)
}
