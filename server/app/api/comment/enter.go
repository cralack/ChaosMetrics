package comment

import (
	"github.com/gin-gonic/gin"
)

type cmntApi struct{}

func InitCommentRouter(router *gin.RouterGroup) {
	api := cmntApi{}
	routerPath := "/comments"
	baseGroup := router.Group(routerPath)
	{
		baseGroup.GET("/list", api.QueryCommentList)
		baseGroup.POST("", api.PostComment)
		baseGroup.DELETE("", api.DeleteComment)
	}
}
