package common

import (
	"github.com/gin-gonic/gin"
)

type cmnApi struct{}

func InitCommonAPI(router *gin.RouterGroup) {
	api := cmnApi{}
	router.GET("/gameversion", api.GetGameVersions)
}
