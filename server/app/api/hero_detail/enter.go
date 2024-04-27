package hero_detail

import (
	"github.com/gin-gonic/gin"
)

type heroDetailApi struct{}

func InitHeroDetailRouter(r *gin.RouterGroup) {
	api := heroDetailApi{}
	routerPath := ""
	routerGroup := r.Group(routerPath)
	{
		routerGroup.GET("/hero_detail", api.QueryHeroDetail)
	}
}
