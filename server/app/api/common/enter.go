package common

import (
	"github.com/gin-gonic/gin"
)

type cmnApi struct{}

func InitCommonAPI(router *gin.RouterGroup) {
	api := cmnApi{}
	router.GET("/gameversion", api.GetGameVersions)
	router.GET("/perks", api.GetPerksData)
	router.GET("/spells", api.GetSpellData)
	router.GET("/items", api.GetItemListData)
}
