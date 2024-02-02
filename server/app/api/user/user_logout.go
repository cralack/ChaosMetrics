package user

import (
	"github.com/cralack/ChaosMetrics/server/app/middleware/auth"
	"github.com/cralack/ChaosMetrics/server/app/provider/user"
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/gin-gonic/gin"
)

// Logout godoc
//
//	@Summary		用户登陆
//	@Description	login @usrname,passwd
//	@Accept			application/json
//	@Produce		application/json
//	@Tags			User Service
//	@Success		200	{object}	response.Response{msg=string}
//	@Router			/user/logout [get]
func (a *usrApi) Logout(ctx *gin.Context) {
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		response.FailWithMessage("user havent login", ctx)
		return
	}

	serv := user.NewUserService()
	if err := serv.Logout(authUser); err != nil {
		response.FailWithMessage("logout failed", ctx)
		return
	}
	response.Ok(ctx)
	return
}
