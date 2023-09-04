package video

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/jinzhu/copier"
	"toktik-api/internal/api"
	"toktik-api/internal/model/request"
	"toktik-api/internal/model/response"
	"toktik-common/errcode"
	res2 "toktik-common/response"
	"toktik-common/token"
	"toktik-rpc/kitex_gen/video"
)

type HandlerVideo struct {
}

func NewHandlerVideo() *HandlerVideo {
	return &HandlerVideo{}
}

func (v *HandlerVideo) VideoFeed(ctx context.Context, c *app.RequestContext) {
	res := res2.NewResponse(c)
	// 1.接收参数 参数模型
	req := &request.VideoFeedRequest{}
	if err := c.Bind(req); err != nil {
		res.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	fmt.Println("req:", req)
	// 若登录，获取用户id
	content, ok := token.GetTokenContent(c)
	if ok {
		req.UserId = content.ID
	}
	// 2.校验参数，这里是校验latest_time
	req.VerifyFeed()
	// 3.调用rpc服务获取响应
	params := &video.VideoFeedRequest{}
	_ = copier.Copy(params, req)
	result, err := api.VideoClient.VideoFeed(ctx, params)
	if err != nil {
		res.Reply(errcode.ErrServer.WithDetails(err.Error()))
		return
	}
	// 4.返回结果
	resp := &response.VideoFeedResponse{}
	_ = copier.Copy(resp, result)
	res.Reply(nil, resp)
}

func (v *HandlerVideo) VideoPublish(ctx context.Context, c *app.RequestContext) {
	res := res2.NewResponse(c)
	// 1.接收参数 参数模型
	req := &request.VideoPublishRequest{}
	if err := c.Bind(req); err != nil {
		res.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	fmt.Println("req:", req)

	// 2.校验参数
	content, ok := token.GetTokenContent(c)
	if !ok {
		res.Reply(errcode.ErrServer)
		return
	}
	req.UserId = content.ID
	dataBuf, err := req.VerifyFeed()
	if err != nil {
		res.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	// 3.调用rpc服务获取响应
	params := &video.VideoPublishRequest{}
	_ = copier.Copy(params, req)
	params.Data = dataBuf.Bytes()
	result, err := api.VideoClient.VideoPublish(ctx, params)
	if err != nil {
		res.Reply(errcode.ErrServer.WithDetails(err.Error()))
		return
	}
	// 4.返回结果
	resp := &response.VideoPublishResponse{}
	_ = copier.Copy(resp, result)
	res.Reply(nil, resp)
}

func (v *HandlerVideo) PublishList(ctx context.Context, c *app.RequestContext) {
	res := res2.NewResponse(c)
	// 1.接收参数 参数模型
	req := &request.PublishListRequest{}
	if err := c.Bind(req); err != nil {
		res.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	fmt.Println("req:", req)

	// 2.校验参数
	content, ok := token.GetTokenContent(c)
	if !ok {
		res.Reply(errcode.ErrServer)
		return
	}
	req.MyUserId = content.ID
	// 3.调用rpc服务获取响应
	params := &video.PublishListRequest{}
	_ = copier.Copy(params, req)
	result, err := api.VideoClient.PublishList(ctx, params)
	if err != nil {
		res.Reply(errcode.ErrServer.WithDetails(err.Error()))
		return
	}
	// 4.返回结果
	resp := &response.PublishListResponse{}
	_ = copier.Copy(resp, result)
	res.Reply(nil, resp)
}
