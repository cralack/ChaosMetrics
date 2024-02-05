package user

import (
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/cralack/ChaosMetrics/server/utils/jwt"
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
	token := utils.GetToken(ctx)
	j := jwt.NewJWT()
	if err := j.SetJWTBlack(token); err != nil {
		response.FailWithDetailed(err, "logout failed", ctx)
		return
	}
	utils.ClearToken(ctx)
	response.Ok(ctx)
	return
}
