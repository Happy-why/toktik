package video

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"toktik-api/internal/api"
	"toktik-api/internal/model/request"
	"toktik-api/internal/model/response"
	"toktik-common/errcode"
	res2 "toktik-common/response"
	"toktik-common/token"
	"toktik-rpc/kitex_gen/comment"
)

type HandlerComment struct {
}

func NewHandlerVideo() *HandlerComment {
	return &HandlerComment{}
}

func (v *HandlerComment) CommentAction(c *gin.Context) {
	res := res2.NewResponse(c)
	// 1.接收参数 参数模型
	req := &request.CommentActionRequest{}
	if err := c.ShouldBind(req); err != nil {
		res.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	// 2.校验参数
	content, ok := token.GetTokenContent(c)
	if !ok {
		res.Reply(errcode.ErrServer)
		return
	}
	req.UserId = content.ID
	// 3.调用rpc服务获取响应
	params := &comment.CommentActionRequest{}
	_ = copier.Copy(params, req)
	result, err := api.Comment.CommentAction(c, params)
	if err != nil {
		res.Reply(errcode.ErrServer.WithDetails(err.Error()))
		return
	}
	// 4.返回结果
	resp := &response.CommentActionResponse{}
	_ = copier.Copy(resp, result)
	res.Reply(nil, resp)
}

func (v *HandlerComment) CommentList(c *gin.Context) {
	res := res2.NewResponse(c)
	// 1.接收参数 参数模型
	req := &request.CommentListRequest{}
	if err := c.ShouldBind(req); err != nil {
		res.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	// 2.校验参数
	// 3.调用rpc服务获取响应
	params := &comment.CommentListRequest{}
	_ = copier.Copy(params, req)
	result, err := api.Comment.CommentList(c, params)
	if err != nil {
		res.Reply(errcode.ErrServer.WithDetails(err.Error()))
		return
	}
	// 4.返回结果
	resp := &response.CommentListResponse{}
	_ = copier.Copy(resp, result)
	res.Reply(nil, resp)
}
