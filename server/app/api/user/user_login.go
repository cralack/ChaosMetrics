package user

import (
	"fmt"
	"time"

	"github.com/cralack/ChaosMetrics/server/app/provider/user"
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model"
	"github.com/cralack/ChaosMetrics/server/model/request"
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/cralack/ChaosMetrics/server/utils/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type loginParam struct {
	UserName string `json:"username" example:"snoop" binding:"required"`
	Password string `json:"password" example:"123456" binding:"required,gte=6"`
}

// Login godoc
//
//	@Summary		用户登陆
//	@Description	login @usrname,passwd
//	@Accept			application/json
//	@Produce		application/json
//	@Tags			User Service
//	@Param			data	body		loginParam	true	"login"
//	@Success		200		{object}	response.Response{msg=string}
//	@Router			/user/login [post]
func (a *usrApi) Login(ctx *gin.Context) {
	var (
		param loginParam
		err   error
		tar   *model.User
		// token string
	)
	// need captha
	if err = ctx.ShouldBindJSON(&param); err != nil {
		response.FailWithMessage("wrong param", ctx)
		return
	}
	if global.ChaEnv != global.ProductEnv {
		time.Sleep(time.Second / 2)
	}
	serv := user.NewUserService()
	if tar, err = serv.Login(param.UserName, param.Password); err != nil {
		response.FailWithMessage(fmt.Sprintf("login failed,%s", err.Error()), ctx)
		return
	}
	a.genToken(ctx, tar)
	return
}

func (a *usrApi) genToken(ctx *gin.Context, tar *model.User) {
	j := jwt.NewJWT()
	claims := j.CreateClaims(request.PrivateClaims{
		UUID:     tar.UUID,
		ID:       tar.ID,
		Username: tar.UserName,
		NickName: tar.NickName,
		Role:     tar.Role,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		global.ChaLogger.Error("get token failed!", zap.Error(err))
		response.FailWithMessage("get token failed", ctx)
		return
	}
	utils.SetToken(ctx, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
	response.OkWithDetailed(gin.H{"token": token}, "login succeed", ctx)
}
