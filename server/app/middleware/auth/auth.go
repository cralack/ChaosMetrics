package auth

import (
	"github.com/cralack/ChaosMetrics/server/app/provider/user"
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/cralack/ChaosMetrics/server/model/usermodel"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get token
		token := utils.GetToken(ctx)
		if token == "" {
			response.FailWithDetailed(gin.H{"reload": true},
				"未登录或非法访问", ctx)
			ctx.Abort()
			return
		}

		// parse token
		serv := user.NewUserService()
		authUser, err := serv.VerifyLogin(token)
		if err != nil {
			response.FailWithDetailed(gin.H{"reload": true},
				"未登录或非法访问", ctx)
			ctx.Abort()
			return
		}
		ctx.Set("auth_user", authUser)
		ctx.Next()
	}
}

func GetAuthUser(ctx *gin.Context) *usermodel.User {
	if tmp, has := ctx.Get("auth_user"); !has {
		return nil
	} else {
		return tmp.(*usermodel.User)
	}
}
