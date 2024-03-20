package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	ERROR   = 4
	QUIET   = 1
	SUCCESS = 0
)

func Result(code int, data interface{}, msg string, c *gin.Context) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func Ok(ctx *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "操作成功", ctx)
}

func OkWithMessage(message string, ctx *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, ctx)
}

func OkWithData(data interface{}, ctx *gin.Context) {
	Result(SUCCESS, data, "查询成功", ctx)
}

func OkWithDetailed(data interface{}, message string, ctx *gin.Context) {
	Result(SUCCESS, data, message, ctx)
}
func OkWithQuiet(data interface{}, ctx *gin.Context) {
	Result(QUIET, data, "查询成功", ctx)
}

func Fail(ctx *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "操作失败", ctx)
}

func FailWithMessage(message string, ctx *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, ctx)
}

func FailWithDetailed(data interface{}, message string, ctx *gin.Context) {
	Result(ERROR, data, message, ctx)
}
