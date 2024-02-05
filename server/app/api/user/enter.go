package user

import (
	"github.com/cralack/ChaosMetrics/server/app/middleware"
	"github.com/gin-gonic/gin"
)

type usrApi struct{}

func InitUserRouter(router *gin.RouterGroup) {
	api := usrApi{}
	routerPath := "/user"
	baseGroup := router.Group(routerPath)
	{
		baseGroup.POST("/register", api.Register)
		baseGroup.GET("/verify", api.Verify)
		baseGroup.POST("/login", api.Login)
	}
	authnGroup := router.Group(routerPath).Use(middleware.JWTAuth())
	{
		authnGroup.GET("/logout", api.Logout)
	}
}
