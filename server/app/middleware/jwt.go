package middleware

import (
	"errors"

	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/cralack/ChaosMetrics/server/utils/jwt"
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := utils.GetToken(ctx)
		if token == "" {
			response.FailWithDetailed(gin.H{"reload": true}, "未登录或非法访问", ctx)
			ctx.Abort()
			return
		}

		j := jwt.NewJWT()
		if j.InBlackList(token) {
			response.FailWithDetailed(gin.H{"reload": true}, "您的帐户异地登陆或令牌失效", ctx)
			utils.ClearToken(ctx)
			ctx.Abort()
			return
		}

		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, jwt.TokenExpired) {
				response.FailWithDetailed(gin.H{"reload": true}, "授权已过期", ctx)
				utils.ClearToken(ctx)
				ctx.Abort()
				return
			}
			response.FailWithDetailed(gin.H{"reload": true}, err.Error(), ctx)
			utils.ClearToken(ctx)
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}
