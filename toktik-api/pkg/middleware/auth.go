package middleware

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"go.uber.org/zap"
	"toktik-api/internal/api"
	"toktik-api/internal/global"
	"toktik-common/errcode"
	"toktik-common/response"
	"toktik-common/token"
	user "toktik-rpc/kitex_gen/user"
)

const (
	TokenNil = ""
	AuthKey  = "token-Content"
)

type AuthToken struct {
	Token string `json:"token" form:"token" query:"token"`
}

// ParseToken 获取并解析header中token
// 返回 payload,token,err
func ParseToken(accessToken string) (*token.Payload, string, errcode.Err) {
	payload, err := global.Maker.VerifyToken(accessToken)
	if err != nil {
		if err.Error() == "超时错误" {
			return nil, "", errcode.ErrAuthExp
		}
		return nil, "", errcode.ErrAuth
	}
	return payload, accessToken, nil
}

// Auth 鉴权中间件,用于解析并写入token
func Auth() func(ctx context.Context, c *app.RequestContext) {
	return func(ctx context.Context, c *app.RequestContext) {
		res := response.NewResponse(c)
		t := &AuthToken{}
		if err := c.Bind(&t); err != nil {
			zap.L().Error("c.ShouldBind(&t) err:", zap.Error(err))
			res.Reply(errcode.ErrAuth.WithDetails(err.Error()))
			c.Abort()
			return
		}
		fmt.Println("token:", t)
		if t.Token == TokenNil {
			c.Next(ctx)
			return
		}
		payload, _, err := ParseToken(t.Token)
		if err != nil {
			c.Next(ctx)
			return
		}
		content := &token.Content{}
		if err := content.Unmarshal(payload.Content); err != nil {
			zap.L().Error("content.Unmarshal(payload.Content) err:", zap.Error(err))
			res.Reply(errcode.ErrAuth.WithDetails(err.Error()))
			c.Abort()
			return
		}
		zap.S().Info(content)
		fmt.Println("content:", content)
		c.Set(AuthKey, content)
		c.Next(ctx)
	}
}

// MustUser 必须是用户
func MustUser() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		res := response.NewResponse(c)
		val, ok := c.Get(AuthKey)
		if !ok {
			res.Reply(errcode.ErrAuth)
			c.Abort()
			return
		}
		content := val.(*token.Content)
		if content.Type != token.UserToken {
			res.Reply(errcode.ErrAuth)
			c.Abort()
			return
		}
		fmt.Println("准备 TokenVerify content:", content)
		_, err := api.UserClient.TokenVerify(ctx, &user.TokenVerifyRequest{UserId: content.ID, TokenType: string(content.Type)})
		if err != nil {
			zap.L().Error("api.UserClient.TokenVerify err:", zap.Error(err))
			res.Reply(errcode.ErrAuth.WithDetails(err.Error()))
			c.Abort()
			return
		}
		fmt.Println("TokenVerify 成功！")
		c.Next(ctx)
	}
}
