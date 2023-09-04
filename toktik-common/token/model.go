package token

import (
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/app"
)

type (
	TokenType string
)

const (
	UserToken TokenType = "user"
	AuthKey             = "token-Content"
)

type Content struct {
	ID   int64     `json:"id"`
	Type TokenType `json:"type"`
}

// NewTokenContent 新建一种类型的token
func NewTokenContent(t TokenType, userID int64) *Content {
	return &Content{Type: t, ID: userID}
}

func (c *Content) Marshal() ([]byte, error) {
	// 返回json编码
	return json.Marshal(c)
}

func (c *Content) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, &c); err != nil {
		return err
	}
	return nil
}

// Token 结合token.Payload和Token
type Token struct {
	AccessToken string
	Payload     *Payload
	Content     *Content
}

// GetTokenContent 从当前上下文中获取保存的content内容
func GetTokenContent(c *app.RequestContext) (*Content, bool) {
	val, ok := c.Get(AuthKey)
	if !ok {
		return nil, false
	}
	return val.(*Content), true
}
