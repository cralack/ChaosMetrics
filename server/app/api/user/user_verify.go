package user

import (
	"github.com/cralack/ChaosMetrics/server/app/provider/user"
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/gin-gonic/gin"
)

type verifyParam struct {
	Token string `form:"token" binding:"required"`
}

// Verify godoc
//
//	@Summary		验证注册信息
//	@Description	query @token
//	@Accept			application/json
//	@Produce		application/json
//	@Tags			User Service
//	@Param			data	query		verifyParam	true	"verify token"
//	@Success		200		{object}	response.Response{msg=string}
//	@Router			/user/verify [get]
func (a *usrApi) Verify(ctx *gin.Context) {
	var (
		param verifyParam
		err   error
		ok    bool
	)
	if err = ctx.ShouldBindQuery(&param); err != nil {
		response.FailWithMessage("wrong param", ctx)
		return
	}

	serv := user.NewUserService()
	ok, err = serv.VerifyRegister(param.Token)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	if !ok {
		response.FailWithMessage("verify failed", ctx)
		return
	}
	response.OkWithMessage("register succeed,refresh", ctx)
}
