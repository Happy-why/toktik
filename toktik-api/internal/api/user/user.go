package user

import (
	"github.com/Happy-Why/toktik-api/internal/api"
	"github.com/Happy-Why/toktik-api/internal/model/request"
	"github.com/Happy-Why/toktik-api/internal/model/response"
	"github.com/Happy-Why/toktik-api/pkg/myerr"
	"github.com/Happy-Why/toktik-common/errcode"
	res2 "github.com/Happy-Why/toktik-common/response"
	"github.com/Happy-Why/toktik-common/token"
	"github.com/Happy-Why/toktik-rpc/kitex_gen/user"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type HandlerUser struct {
}

func NewHandlerUser() *HandlerUser {
	return &HandlerUser{}
}

func (h *HandlerUser) Register(c *gin.Context) {
	res := res2.NewResponse(c)
	// 1.接收参数 参数模型
	req := &request.RegisterRequest{}
	if err := c.ShouldBind(req); err != nil {
		res.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	// 2.校验参数 判断参数是否合法
	if err := req.Verify(); err != nil {
		res.Reply(myerr.ErrUserNameORPassWord)
		return
	}
	// 3.调用user rpc服务 获取响应
	params := &user.RegisterRequest{}
	_ = copier.Copy(params, req)
	result, err := api.UserClient.Register(c, params)
	if err != nil {
		res.Reply(errcode.ErrServer.WithDetails(err.Error()))
		return
	}
	// 4.返回结果
	resp := &response.RegisterResponse{}
	_ = copier.Copy(resp, result)
	res.Reply(nil, resp)
	return
}

func (h *HandlerUser) Login(c *gin.Context) {
	res := res2.NewResponse(c)
	// 1.接收参数 参数模型
	req := &request.LoginRequest{}
	if err := c.ShouldBind(req); err != nil {
		res.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	// 2.校验参数 判断参数是否合法
	if err := req.Verify(); err != nil {
		res.Reply(myerr.ErrUserNameORPassWord)
		return
	}
	// 3.调用user rpc服务 获取响应
	params := &user.LoginRequest{}
	_ = copier.Copy(params, req)
	resp, err := api.UserClient.Login(c, params)
	if err != nil {
		res.Reply(errcode.ErrServer.WithDetails(err.Error()))
		return
	}
	// 4.返回结果
	res.Reply(nil, resp)
	return
}

func (h *HandlerUser) UserIndex(c *gin.Context) {
	res := res2.NewResponse(c)
	// 1.接收参数 参数模型
	req := &request.UserIndexRequest{}
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
	// 2.调用user rpc服务 获取响应
	params := &user.UserIndexRequest{}
	_ = copier.Copy(params, req)

	result, err := api.UserClient.UserIndex(c, params)
	if err != nil {
		res.Reply(errcode.ErrServer.WithDetails(err.Error()))
		return
	}
	// 4.返回结果
	resp := &response.UserIndexResponse{}
	_ = copier.Copy(resp, result)
	res.Reply(nil, resp)
	return
}
