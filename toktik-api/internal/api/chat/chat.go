package chat

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"toktik-api/internal/api"
	"toktik-api/internal/model/request"
	"toktik-api/internal/model/response"
	"toktik-api/pkg/myerr"
	"toktik-common/errcode"
	res2 "toktik-common/response"
	"toktik-common/token"
	"toktik-rpc/kitex_gen/chat"
)

type HandlerChat struct {
}

func NewHandlerInteraction() *HandlerChat {
	return &HandlerChat{}
}

func (h *HandlerChat) ChatAction(c *gin.Context) {
	res := res2.NewResponse(c)
	// 1.接收参数 参数模型
	req := &request.ChatActionRequest{}
	if err := c.ShouldBind(req); err != nil {
		res.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := token.GetTokenContent(c)
	if !ok {
		res.Reply(errcode.ErrServer)
		return
	}
	req.UserId = content.ID
	if req.Verify() {
		res.Reply(myerr.CanNotChatSelf)
	}
	// 2.调用 rpc服务 获取响应
	params := &chat.ChatActionRequest{}
	_ = copier.Copy(params, req)
	result, err := api.ChatClient.ChatAction(c, params)
	if err != nil {
		res.Reply(errcode.ErrServer.WithDetails(err.Error()))
		return
	}
	// 3.返回结果
	resp := &response.ChatActionResponse{}
	_ = copier.Copy(resp, result)
	res.Reply(nil, resp)
}

func (h *HandlerChat) MessageList(c *gin.Context) {
	res := res2.NewResponse(c)
	// 1.接收参数 参数模型
	req := &request.MessageListRequest{}
	if err := c.ShouldBind(req); err != nil {
		res.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := token.GetTokenContent(c)
	if !ok {
		res.Reply(errcode.ErrServer)
		return
	}
	req.UserId = content.ID

	// 2.调用 rpc服务 获取响应
	params := &chat.MessageListRequest{}
	_ = copier.Copy(params, req)
	result, err := api.ChatClient.MessageList(c, params)
	if err != nil {
		res.Reply(errcode.ErrServer.WithDetails(err.Error()))
		return
	}
	// 3.返回结果
	resp := &response.MessageListResponse{}
	_ = copier.Copy(resp, result)
	res.Reply(nil, resp)
}
