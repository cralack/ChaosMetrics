package item

import (
	"github.com/cralack/ChaosMetrics/server/internal/service/provider/item"
	"github.com/gin-gonic/gin"
)

type itemQueryParam struct {
	ItemID  string `json:"itemid" binding:"required"`
	Lang    string `json:"lang" binding:"required"`
	Version string `json:"version" binding:"required"`
}

func (i *itemApi) QueryApi(ctx *gin.Context) {
	itemService := item.NewItemService()
	var param itemQueryParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSON(400, gin.H{
			"msg": "wrong param",
		})
		return
	}

	if itemDTO, err := itemService.QueryItem(param.ItemID, param.Version, param.Lang); err != nil {
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
