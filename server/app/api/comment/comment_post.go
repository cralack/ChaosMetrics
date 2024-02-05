package comment

import (
	"github.com/cralack/ChaosMetrics/server/app/provider/comment"
	"github.com/cralack/ChaosMetrics/server/model"
	"github.com/cralack/ChaosMetrics/server/model/response"
	"github.com/cralack/ChaosMetrics/server/utils"
	"github.com/gin-gonic/gin"
)

type postCommentParam struct {
	ChampionID string `json:"championID" example:"Ahri" binding:"required"` // Champion name
	Version    string `json:"version" example:"14.1.1" binding:"required"`  // Version
	Comment    string `json:"comment" example:"so cute" binding:"required"`
}

// PostComment godoc
//
//	@Summary		发表一个评论
//	@Description	post @ChampionID,Version
//	@Accept			application/json
//	@Produce		application/json
//	@Tags			Comment
//	@Param			postCommentParam	body		postCommentParam	true	"Post a comment @ champion,version"
//	@Success		200					{object}	response.Response{msg=string}
//	@Router			/comments [post]
func (a *cmntApi) PostComment(ctx *gin.Context) {
	var (
		param postCommentParam
		err   error
	)
	if err = ctx.ShouldBindJSON(&param); err != nil {
		response.FailWithMessage("wrong param", ctx)
		return
	}
	id := utils.GetUserID(ctx)
	serv := comment.NewCommentService()
	if err = serv.PostComment(&model.Comment{
		ChampionID: param.ChampionID,
		Version:    param.Version,
		Content:    param.Comment,
		AuthorID:   id,
	}); err != nil {
		response.FailWithDetailed(err, "post comment failed", ctx)
		return
	}
	response.Ok(ctx)
}
