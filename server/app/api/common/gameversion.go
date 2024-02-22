package common

import (
	"context"
	"encoding/json"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/gin-gonic/gin"
)

func (a *cmnApi) GetGameVersions(ctx *gin.Context) {
	res := global.ChaRDB.HGet(context.Background(), "/version", "versions")
	if res.Err() != nil {
		response.Fail(ctx)
		return
	}
	versions := make([]string, 0)
	if err := json.Unmarshal([]byte(res.Val()), &versions); err != nil {
		response.Fail(ctx)
		return
	}
	idx := 0
	for i, v := range versions {
		if v == "13.1.1" {
			idx = i + 1
			break
		}
	}
	versions = versions[:idx]
	response.OkWithQuiet(versions, ctx)
}
