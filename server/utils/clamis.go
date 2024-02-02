package utils

import (
	"net"

	"github.com/gin-gonic/gin"
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
