package user

import (
	"time"

	"github.com/cralack/ChaosMetrics/server/app/provider/user"
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model/request"
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/cralack/ChaosMetrics/server/model/usermodel"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type loginParam struct {
	UserName string `json:"username" example:"snoop" binding:"required"`
	Password string `json:"password" example:"123456" binding:"required,gte=6"`
}
type LoginResponse struct {
	User      usermodel.User `json:"user"`
	Token     string         `json:"token"`
	ExpiresAt int64          `json:"expiresAt"`
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
		token string
	)
	// need captha
	if err = ctx.ShouldBindJSON(&param); err != nil {
		response.FailWithMessage("wrong param", ctx)
	}
	serv := user.NewUserService()
	if token, err = serv.Login(param.UserName, param.Password); err != nil {
		response.FailWithDetailed(err, "login failed", ctx)
	}
	utils.SetToken(ctx, token, 3600)
	// a.TokenNext(ctx)
	response.Ok(ctx)
	return
}

func (a *usrApi) TokenNext(ctx *gin.Context, tar usermodel.User) {
	j := utils.NewJWT()
	claims := j.CreateClaims(request.BaseClaims{
		UUID:        tar.UUID,
		ID:          tar.ID,
		Username:    tar.UserName,
		NickName:    tar.NickName,
		AuthorityId: tar.Role,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		global.ChaLogger.Error("get token failed!", zap.Error(err))
		response.FailWithMessage("get token failed", ctx)
		return
	}
	utils.SetToken(ctx, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
}
