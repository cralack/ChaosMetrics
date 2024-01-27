package item

import (
	"github.com/gin-gonic/gin"
)

type itemApi struct{}

func InitItemRouter(r *gin.Engine) {
	api := itemApi{}
	itemApi := r.Group("/items")
	{
		itemApi.GET("/item", api.QueryApi)
	}
}
