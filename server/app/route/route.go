package route

import (
	"net/http"

	"github.com/cralack/ChaosMetrics/server/app/api/item"
	"github.com/cralack/ChaosMetrics/server/internal/global"
	"github.com/gin-gonic/gin"
)

func RegiserRoutes(r *gin.Engine) {
	PublicGroup := r.Group("")
	{
		PublicGroup.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, "ok")
		})
		item.InitItemRouter(r)
	}
	global.ChaLogger.Info("router regist success")
}
