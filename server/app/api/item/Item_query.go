package item

import (
	"github.com/cralack/ChaosMetrics/server/internal/service/provider/item"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/gin-gonic/gin"
)

// itemQueryParam represents the query parameters for an item request
// @Description Query parameters for requesting item details
type itemQueryParam struct {
	ItemID  string `form:"itemid" default:"2010" binding:"required"`    // The ID of the item
	Lang    string `form:"lang" default:"zh_CN" binding:"required"`     // Language
	Version string `form:"version" default:"13.8.1" binding:"required"` // Version
}

// QueryApi godoc
//
//	@Summary		请求一个物品详情
//	@Description	请求一个物品详情 @version,lang
//	@Accept			application/json
//	@Produce		application/json
//	@Tags			item
//	@Param			itemQueryParam	query		itemQueryParam	true	"Query parameters for item"
//	@Success		200				{object}	riotmodel.ItemDTO
//	@Router			/items/item [get]
func (i *itemApi) QueryApi(ctx *gin.Context) {
	var (
		param   itemQueryParam
		itemDTO *riotmodel.ItemDTO
		err     error
	)
	itemService := item.NewItemService()

	if err := ctx.ShouldBindQuery(&param); err != nil {
		ctx.JSON(400, gin.H{
			"msg": "wrong param",
		})
		return
	}

	if itemDTO, err = itemService.QueryItem(param.ItemID, param.Version, param.Lang); err != nil {
		ctx.JSON(400, gin.H{
			"msg": "can not find item by " + err.Error(),
		})
		return
	} else {
		ctx.JSON(200, gin.H{
			"msg":  "query success",
			"data": itemDTO,
		})
	}
}
