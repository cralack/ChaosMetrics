package item

import (
	"github.com/cralack/ChaosMetrics/server/app/provider/item"
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/gin-gonic/gin"
)

// itemQueryParam represents the query parameters for an item request
//
//	@Description	Query parameters for requesting item details
type itemQueryParam struct {
	ItemID  string `form:"itemid" default:"2010" binding:"required"`    // The ID of the item
	Lang    string `form:"lang" default:"zh_CN" binding:"required"`     // Language
	Version string `form:"version" default:"13.8.1" binding:"required"` // Version
}

// QueryApi godoc
//
//	@Summary		请求一个物品详情
//	@Description	query @version,lang
//	@Accept			application/json
//	@Produce		application/json
//	@Tags			Item
//	@Param			itemQueryParam	query		itemQueryParam	true	"Query parameters for item"
//	@Success		200				{object}	response.Response{data=riotmodel.ItemDTO}
//	@Router			/item [get]
func (i *itemApi) QueryApi(ctx *gin.Context) {
	var (
		param   itemQueryParam
		itemDTO *riotmodel.ItemDTO
		err     error
	)
	itemService := item.NewItemService()

	if err = ctx.ShouldBindQuery(&param); err != nil {
		response.FailWithMessage("wrong param", ctx)
		return
	}

	if itemDTO, err = itemService.QueryItem(param.ItemID, param.Version, param.Lang); err != nil {
		response.FailWithDetailed(err, "can not find item", ctx)
		return
	} else {
		response.OkWithData(itemDTO, ctx)
	}
}
