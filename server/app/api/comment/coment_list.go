package comment

import (
	"github.com/cralack/ChaosMetrics/server/app/provider/comment"
	"github.com/cralack/ChaosMetrics/server/model"
	"github.com/cralack/ChaosMetrics/server/model/request"
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/gin-gonic/gin"
)

type queryCommentParam struct {
	ChampionID string `form:"championID" default:"Ahri" binding:"required"` // Champion name
	Version    string `form:"version" default:"14.1.1" binding:"required"`  // Version
	Start      int    `form:"start" default:"0"`
	Size       int    `form:"size" default:"50" binding:"required"`
}

// QueryCommentList godoc
//
//	@Summary		请求一个英雄的评论
//	@Description	query @ChampionID,Version,start,size
//	@Accept			application/json
//	@Produce		application/json
//	@Tags			Comment
//	@Param			squeryCommentParam	query		queryCommentParam	true	"query a comments list"
//	@Success		200					{object}	response.Response{data=[]CommentsDTO,msg=string}
//	@Router			/comments/list [get]
func (a *cmntApi) QueryCommentList(ctx *gin.Context) {
	var (
		param    queryCommentParam
		comments []*model.Comment
		res      []*CommentsDTO
		err      error
	)
	if err = ctx.ShouldBindQuery(&param); err != nil {
		response.FailWithMessage("wrong param", ctx)
		return
	}
	serv := comment.NewCommentService()

	pager := &request.Pager{
		Start: param.Start,
		Size:  param.Size,
	}
	comments, err = serv.GetComments(param.ChampionID, param.Version, pager)
	if err != nil {
		response.FailWithDetailed(err, "get comments failed", ctx)
		return
	}
	res = ConvertCommentsDTO(comments...)
	response.OkWithData(gin.H{
		"pager":    pager,
		"comments": res,
	}, ctx)
}
