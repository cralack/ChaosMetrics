package item

import (
	"github.com/gin-gonic/gin"
)

type itemApi struct{}

func InitItemRouter(router *gin.RouterGroup) {
	api := itemApi{}
	routerPath := ""
	routerGroup := router.Group(routerPath)
	{
		routerGroup.GET("/item", api.QueryApi)
	}
}
