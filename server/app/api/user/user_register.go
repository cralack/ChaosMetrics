package user

import (
	"github.com/cralack/ChaosMetrics/server/app/provider/user"
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/cralack/ChaosMetrics/server/model/usermodel"
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
		tar   *usermodel.User
		token string
		err   error
	)
	if err = ctx.ShouldBindJSON(&param); err != nil {
		response.FailWithMessage("wrong param", ctx)
		return
	}

	serv := user.NewUserService()
	tar = &usermodel.User{
		UserName: param.UserName,
		Password: param.Password,
		Email:    param.Email,
	}
	if token, err = serv.PreRegister(tar); err != nil {
		response.FailWithDetailed(err, "register failed,try later", ctx)
	}
	global.ChaLogger.Debug(token)
	if err = serv.SendVerifyEmail(token); err != nil {
		response.FailWithMessage("send mail ailed,try later", ctx)
	}
	response.Ok(ctx)
	return
}
