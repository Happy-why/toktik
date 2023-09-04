package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"toktik-api/internal/api"
	"toktik-api/internal/model/request"
	"toktik-api/internal/model/response"
	"toktik-api/pkg/myerr"
	"toktik-common/errcode"
	res2 "toktik-common/response"
	"toktik-common/token"
	"toktik-rpc/kitex_gen/user"
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
	result, err := api.UserClient.Login(c, params)
	if err != nil {
		res.Reply(errcode.ErrServer.WithDetails(err.Error()))
		return
	}
	fmt.Println("result:", result)
	// 4.返回结果
	resp := &response.LoginResponse{}
	_ = copier.Copy(resp, result)
	fmt.Println("resp:", resp)
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
