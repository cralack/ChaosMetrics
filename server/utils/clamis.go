package utils

import (
	"net"

	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/cralack/ChaosMetrics/server/model/request"
	"github.com/cralack/ChaosMetrics/server/utils/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ClearToken(ctx *gin.Context) {
	host, _, err := net.SplitHostPort(ctx.Request.Host)
	if err != nil {
		host = ctx.Request.Host
	}
	ctx.SetCookie(global.TokenKey, "", -1,
		"/", host, true, false)

}

func SetToken(ctx *gin.Context, token string, maxAge int) {
	host, _, err := net.SplitHostPort(ctx.Request.Host)
	if err != nil {
		host = ctx.Request.Host
	}
	ctx.SetCookie(global.TokenKey, token, maxAge,
		"/", host, true, false)
}

func GetToken(ctx *gin.Context) string {
	token, _ := ctx.Cookie(global.TokenKey)
	if token == "" {
		token = ctx.Request.Header.Get(global.TokenKey)
	}
	return token
}

func GetClaims(ctx *gin.Context) (*request.CustomClaims, error) {
	token := GetToken(ctx)
	j := jwt.NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		global.ChaLogger.Error("从Gin的Context中获取从jwt解析信息失败, 请检查请求头是否存在x-token且claims是否为规定结构")
	}
	return claims, err
}

// GetUserUuid 从Gin的Context中获取从jwt解析出来的用户UUID
func GetUserUuid(ctx *gin.Context) uuid.UUID {
	if claims, exists := ctx.Get("claims"); exists {
		if waitUse, ok := claims.(*request.CustomClaims); ok {
			return waitUse.UUID
		}
		// 类型断言失败时处理
		return uuid.UUID{}
	}

	cl, err := GetClaims(ctx)
	if err != nil || cl == nil {
		return uuid.UUID{}
	}
	return cl.UUID
}

// GetUserID 从Gin的Context中获取从jwt解析出来的用户ID
func GetUserID(ctx *gin.Context) uint {
	if claims, exists := ctx.Get("claims"); exists {
		if waitUse, ok := claims.(*request.CustomClaims); ok {
			return waitUse.PrivateClaims.ID
		}
		// 类型断言失败时处理
		return 0
	}

	cl, err := GetClaims(ctx)
	if err != nil || cl == nil {
		return 0
	}
	return cl.PrivateClaims.ID
}
