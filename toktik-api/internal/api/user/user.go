package user

import (
	"github.com/Happy-Why/toktik-api/internal/model/request"
	"github.com/Happy-Why/toktik-api/pkg/myerr"
	"github.com/Happy-Why/toktik-common/errcode"
	"github.com/Happy-Why/toktik-common/response"
	"github.com/Happy-Why/toktik-rpc/kitex_gen/user"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"time"
)

type HandlerUser struct {
}

func NewHandlerUser() *HandlerUser {
	return &HandlerUser{}
}

func (h *HandlerUser) Register(c *gin.Context) {
	res := response.NewResponse(c)
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
	resp, err := userClient.Register(c, params, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		res.Reply(errcode.ErrServer.WithDetails(err.Error()))
		return
	}
	// 4.返回结果
	res.Reply(nil, resp)
	return
}
