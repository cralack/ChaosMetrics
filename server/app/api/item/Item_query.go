package item

import (
	"github.com/cralack/ChaosMetrics/server/internal/service/provider/item"
	"github.com/cralack/ChaosMetrics/server/model/riotmodel"
	"github.com/gin-gonic/gin"
)

type itemQueryParam struct {
	ItemID  string `json:"itemid" binding:"required"`
	Lang    string `json:"lang" binding:"required"`
	Version string `json:"version" binding:"required"`
}

// QueryApi godoc
// @Summary		获得item
// @Description	获得item@ersion in lang详情
// @Accept			application/json
// @Produce		application/json
// @Tags			item
// @Param			itemQueryParam	body		itemQueryParam	true	"query with param"
// @Success		200				{object}	riotmodel.ItemDTO
// @Router			/items/item [post]
func (i *itemApi) QueryApi(ctx *gin.Context) {
	var (
		param   itemQueryParam
		itemDTO *riotmodel.ItemDTO
		err     error
	)
	itemService := item.NewItemService()

	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSON(400, gin.H{
			"msg": "wrong param",
		})
		return
	}

	if itemDTO, err = itemService.QueryItem(param.ItemID, param.Version, param.Lang); err != nil {
		ctx.JSON(400, gin.H{
			"msg": "can not find item by" + err.Error(),
		})
		return
	} else {
		ctx.JSON(200, gin.H{
			"msg":  "query success",
			"data": itemDTO,
		})
	}
}
