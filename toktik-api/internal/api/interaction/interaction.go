package interaction

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"toktik-api/internal/api"
	"toktik-api/internal/model/request"
	"toktik-api/internal/model/response"
	"toktik-common/errcode"
	res2 "toktik-common/response"
	"toktik-common/token"
	inter "toktik-rpc/kitex_gen/interaction"
)

type HandlerInteraction struct {
}

func NewHandlerInteraction() *HandlerInteraction {
	return &HandlerInteraction{}
}

func (p *HandlerInteraction) FollowSB(c *gin.Context) {
	res := res2.NewResponse(c)
	// 1.接收参数 参数模型
	req := &request.FollowActionRequest{}
	if err := c.ShouldBind(req); err != nil {
		res.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := token.GetTokenContent(c)
	if !ok {
		res.Reply(errcode.ErrServer)
		return
	}
	req.MyUserId = content.ID

	// 2.调用 rpc服务 获取响应
	params := &inter.FollowActionRequest{}
	_ = copier.Copy(params, req)
	result, err := api.InteractionClient.FollowSB(c, params)
	if err != nil {
		res.Reply(errcode.ErrServer.WithDetails(err.Error()))
		return
	}
	// 3.返回结果
	resp := &response.FollowActionResponse{}
	_ = copier.Copy(resp, result)
	res.Reply(nil, resp)
	return
}

func (p *HandlerInteraction) FollowList(c *gin.Context) {
	res := res2.NewResponse(c)
	// 1.接收参数 参数模型
	req := &request.FollowListRequest{}
	if err := c.ShouldBind(req); err != nil {
		res.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	// 2.调用 rpc服务 获取响应
	params := &inter.FollowListRequest{}
	_ = copier.Copy(params, req)
	result, err := api.InteractionClient.FollowList(c, params)
	if err != nil {
		res.Reply(errcode.ErrServer.WithDetails(err.Error()))
		return
	}
	// 3.返回结果
	resp := &response.FollowListResponse{}
	_ = copier.Copy(resp, result)
	res.Reply(nil, resp)
	return
}

func (p *HandlerInteraction) FansList(c *gin.Context) {
	res := res2.NewResponse(c)
	// 1.接收参数 参数模型
	req := &request.FansListRequest{}
	if err := c.ShouldBind(req); err != nil {
		res.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	// 2.调用 rpc服务 获取响应
	params := &inter.FansListRequest{}
	_ = copier.Copy(params, req)
	result, err := api.InteractionClient.FansList(c, params)
	if err != nil {
		res.Reply(errcode.ErrServer.WithDetails(err.Error()))
		return
	}
	// 3.返回结果
	resp := &response.FansListResponse{}
	_ = copier.Copy(resp, result)
	res.Reply(nil, resp)
	return
}

func (p *HandlerInteraction) FriendList(c *gin.Context) {
	res := res2.NewResponse(c)
	// 1.接收参数 参数模型
	req := &request.FriendListRequest{}
	if err := c.ShouldBind(req); err != nil {
		res.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	// 2.调用 rpc服务 获取响应
	params := &inter.FriendListRequest{}
	_ = copier.Copy(params, req)
	result, err := api.InteractionClient.FriendList(c, params)
	if err != nil {
		res.Reply(errcode.ErrServer.WithDetails(err.Error()))
		return
	}
	// 3.返回结果
	resp := &response.FriendListResponse{}
	_ = copier.Copy(resp, result)
	res.Reply(nil, resp)
	return
}
