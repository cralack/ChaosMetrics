package user

import (
	"github.com/cralack/ChaosMetrics/server/app/provider/user"
	"github.com/cralack/ChaosMetrics/server/model"
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetUserInfo godoc
//
//	@Summary	获取用户信息
//	@Description
//	@Accept		application/json
//	@Produce	application/json
//	@Tags		User Service
//	@Security	TokenAuth
//	@Success	200	{object}	response.Response{data=map[string]interface{},msg=string}
//	@Router		/user/info [get]
func (a *usrApi) GetUserInfo(ctx *gin.Context) {
	var (
		// param queryUserinfoParam
		uid uuid.UUID
		err error
		tar *model.User
		// token string
	)

	uid = utils.GetUserUuid(ctx)
	serv := user.NewUserService()
	if tar, err = serv.GetUserIno(uid); err != nil {
		response.FailWithMessage("get user info failed", ctx)
		return
	}
	response.OkWithQuiet(tar, ctx)
}
