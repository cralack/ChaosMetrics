package common

import (
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/gin-gonic/gin"
)

func (a *cmnApi) GetGameVersions(ctx *gin.Context) {
	res := utils.GetCurMajorVersions()
	if len(res) == 0 {
		response.Fail(ctx)
		return
	}

	response.OkWithQuiet(res, ctx)
}
