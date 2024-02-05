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
	ctx.SetCookie("x-token", "", -1,
		"/", host, true, false)

}

func SetToken(ctx *gin.Context, token string, maxAge int) {
	host, _, err := net.SplitHostPort(ctx.Request.Host)
	if err != nil {
		host = ctx.Request.Host
	}
	ctx.SetCookie("x-token", token, maxAge,
		"/", host, true, false)
}

func GetToken(ctx *gin.Context) string {
	token, _ := ctx.Cookie("x-token")
	if token == "" {
		token = ctx.Request.Header.Get("x-token")
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
	if claims, exists := ctx.Get("claims"); !exists {
		if cl, err := GetClaims(ctx); err != nil {
			return uuid.UUID{}
		} else {
			return cl.UUID
		}
	} else {
		waitUse := claims.(*request.CustomClaims)
		return waitUse.UUID
	}
}

// GetUserID 从Gin的Context中获取从jwt解析出来的用户ID
func GetUserID(ctx *gin.Context) uint {
	if claims, exists := ctx.Get("claims"); !exists {
		if cl, err := GetClaims(ctx); err != nil {
			return 0
		} else {
			return cl.PrivateClaims.ID
		}
	} else {
		waitUse := claims.(*request.CustomClaims)
		return waitUse.PrivateClaims.ID
	}
}
