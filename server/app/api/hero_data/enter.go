package hero_data

import (
	"github.com/gin-gonic/gin"
)

type heroDataApi struct{}

func InitHeroDataRouter(r *gin.RouterGroup) {
	api := heroDataApi{}
	routerPath := ""
	routerGroup := r.Group(routerPath)
	{
		routerGroup.GET("/hero", api.QueryHeroData)
	}
}
