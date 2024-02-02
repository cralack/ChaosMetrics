package item

import (
	"github.com/gin-gonic/gin"
)

type itemApi struct{}

func InitItemRouter(router *gin.RouterGroup) {
	api := itemApi{}
	routerPath := ""
	itemRouter := router.Group(routerPath)
	{
		itemRouter.GET("/item", api.QueryApi)
	}
}
