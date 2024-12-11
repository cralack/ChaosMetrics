package middleware

import (
	"strconv"
	"strings"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/cralack/ChaosMetrics/server/utils/casbin"
	"github.com/gin-gonic/gin"
)

var casbinServ = casbin.CbServ

func CasbinHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if global.ChaEnv == global.ProductEnv {
			waitUse, _ := utils.GetClaims(ctx)
			// 获取请求的PATH
			path := ctx.Request.URL.Path
			obj := strings.TrimPrefix(path, global.ChaConf.Router.RouterPrefix)
			// 获取请求方法
			act := ctx.Request.Method
			// 获取用户的角色
			sub := strconv.Itoa(int(waitUse.Role))
			e := casbinServ.Casbin()
			success, _ := e.Enforce(sub, obj, act)
			if !success {
				response.FailWithMessage("action forbidn", ctx)
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}
