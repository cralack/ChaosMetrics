package common

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/cralack/ChaosMetrics/server/utils"
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
	majorVersion, _ := strconv.Atoi(versions[0][:2])
	majorVersion *= 100
	for i, v := range versions {
		if ver, _ := utils.ConvertVersionToIdx(v); int(ver) <= majorVersion {
			versions = versions[:i]
			response.OkWithQuiet(versions, ctx)
			return
		}
	}
}
