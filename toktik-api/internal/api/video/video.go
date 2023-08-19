package video

import (
	"github.com/Happy-Why/toktik-api/internal/model/request"
	"github.com/Happy-Why/toktik-common/errcode"
	res2 "github.com/Happy-Why/toktik-common/response"
	"github.com/gin-gonic/gin"
)

type HandlerVideo struct {
}

func NewHandlerVideo() *HandlerVideo {
	return &HandlerVideo{}
}

func (v *HandlerVideo) VideoFeed(c *gin.Context) {
	res := res2.NewResponse(c)
	// 1.接收参数 参数模型
	req := &request.FollowActionRequest{}
	if err := c.ShouldBind(req); err != nil {
		res.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
}
