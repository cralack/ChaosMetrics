package comment

import (
	"github.com/cralack/ChaosMetrics/server/app/provider/comment"
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/gin-gonic/gin"
)

type delCommentParam struct {
	ID uint `form:"id"  binding:"required"`
}

// DeleteComment godoc
//
//	@Summary		删除一个评论
//	@Description	delete @id
//	@Accept			application/json
//	@Produce		application/json
//	@Tags			Comment
//	@Param			delCommentParam	query		delCommentParam	true	"delete a comment"
//	@Success		200				{object}	response.Response{msg=string}
//	@Router			/comments [delete]
func (a *cmntApi) DeleteComment(ctx *gin.Context) {
	var (
		param delCommentParam
		err   error
	)
	if err = ctx.ShouldBindQuery(&param); err != nil {
		response.FailWithMessage("wrong param", ctx)
		return
	}
	serv := comment.NewCommentService()
	if err = serv.DeleteComments(param.ID); err != nil {
		response.FailWithDetailed(err, "delete failed", ctx)
		return
	}
	response.Ok(ctx)
}
