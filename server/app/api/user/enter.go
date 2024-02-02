package user

import (
	"github.com/cralack/ChaosMetrics/server/app/middleware/auth"
	"github.com/gin-gonic/gin"
)

type usrApi struct{}

func InitUserRouter(router *gin.RouterGroup) {
	api := usrApi{}
	routerPath := "/user"
	userRouter := router.Group(routerPath).Use()
	{
		userRouter.POST("/register", api.Register)
		userRouter.GET("/verify", api.Verify)
		userRouter.POST("/login", api.Login)
		userRouter.GET("/logout", api.Logout, auth.Auth())
	}

}
