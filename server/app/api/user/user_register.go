package user

import (
	"github.com/cralack/ChaosMetrics/server/app/provider/user"
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model"
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/gin-gonic/gin"
)

type registerParam struct {
	UserName string `json:"username" example:"snoop" binding:"required"`
	Password string `json:"password" example:"123456" binding:"required,gte=6"`
	Email    string `json:"email" example:"snoop@dogg.com" binding:"required,gte=6"`
}

// Register godoc
//
//	@Summary		注册一个新用户
//	@Description	register @usrname,passwd,email
//	@Accept			application/json
//	@Produce		application/json
//	@Tags			User Service
//	@Param			data	body		registerParam	true	"Register a new user"
//	@Success		200		{object}	response.Response{msg=string}
//	@Router			/user/register [post]
func (a *usrApi) Register(ctx *gin.Context) {
	var (
		param registerParam
		tar   *model.User
		token string
		err   error
	)
	if err = ctx.ShouldBindJSON(&param); err != nil {
		response.FailWithMessage("wrong param", ctx)
		return
	}

	serv := user.NewUserService()
	tar = &model.User{
		UserName: param.UserName,
		Password: param.Password,
		Email:    param.Email,
	}
	if token, err = serv.PreRegister(tar); err != nil {
		response.FailWithDetailed(err, "register failed,try later", ctx)
		return
	}
	global.ChaLogger.Debug(token)
	if err = serv.SendVerifyEmail(tar, token); err != nil {
		response.FailWithMessage("send mail failed,try later", ctx)
		return
	}
	response.Ok(ctx)
	return
}
