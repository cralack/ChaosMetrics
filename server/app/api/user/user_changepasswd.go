package user

import (
	"github.com/cralack/ChaosMetrics/server/app/provider/user"
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/gin-gonic/gin"
)

type changePasswordParam struct {
	Password    string `json:"password" binding:"required"`    // 密码
	NewPassword string `json:"newPassword" binding:"required"` // 新密码
}

// ChangePasswd godoc
//
//	@Summary		更改密码
//	@Description	change @passwd,newpasswd
//	@Accept			application/json
//	@Produce		application/json
//	@Security		TokenAuth
//	@Tags			User Service
//	@Param			data	body		changePasswordParam	true	"login"
//	@Success		200		{object}	response.Response{msg=string}
//	@Router			/user/changepasswd [post]
func (a *usrApi) ChangePasswd(ctx *gin.Context) {
	var (
		param changePasswordParam
		err   error
		id    uint
	)
	if err = ctx.ShouldBindJSON(&param); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	serv := user.NewUserService()
	id = utils.GetUserID(ctx)
	if err = serv.ChangePassword(id, param.Password, param.NewPassword); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}
	response.Ok(ctx)
}
