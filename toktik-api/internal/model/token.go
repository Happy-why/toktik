package model

import (
	"encoding/json"
	"github.com/Happy-Why/toktik-common/token"
)

type TokenType string

const (
	UserToken    TokenType = "user"
	AccountToken TokenType = "account"
)

type Content struct {
	Type TokenType `json:"type"`
	ID   int64     `json:"id"`
}

// NewTokenContent 新建一种类型的token
func NewTokenContent(t TokenType, ID int64) *Content {
	return &Content{Type: t, ID: ID}
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
	Payload     *token.Payload
	Content     *Content
}
