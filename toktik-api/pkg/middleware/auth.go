package middleware

import (
	"github.com/Happy-Why/toktik-api/internal/api"
	"github.com/Happy-Why/toktik-api/internal/global"
	"github.com/Happy-Why/toktik-common/errcode"
	"github.com/Happy-Why/toktik-common/response"
	"github.com/Happy-Why/toktik-common/token"
	user "github.com/Happy-Why/toktik-rpc/kitex_gen/user"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	TokenNil = ""
	AuthKey  = "token-Content"
)

type AuthToken struct {
	Token string `json:"token" form:"token"`
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
func Auth() func(c *gin.Context) {
	return func(c *gin.Context) {
		res := response.NewResponse(c)
		t := &AuthToken{}
		if err := c.ShouldBind(&t); err != nil {
			zap.L().Error("c.ShouldBind(&t) err:", zap.Error(err))
			res.Reply(errcode.ErrAuth.WithDetails(err.Error()))
			c.Abort()
			return
		}
		if t.Token == TokenNil {
			c.Next()
			return
		}
		payload, _, err := ParseToken(t.Token)
		if err != nil {
			c.Next()
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
		c.Set(AuthKey, content)
		c.Next()
	}
}

// MustUser 必须是用户
func MustUser() gin.HandlerFunc {
	return func(c *gin.Context) {
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
		_, err := api.UserClient.TokenVerify(c,
			&user.TokenVerifyRequest{UserId: content.ID, TokenType: string(content.Type)},
		)
		if err != nil {
			zap.L().Error("api.UserClient.TokenVerify err:", zap.Error(err))
			res.Reply(errcode.ErrAuth.WithDetails(err.Error()))
			c.Abort()
			return
		}
		c.Next()
	}
}
