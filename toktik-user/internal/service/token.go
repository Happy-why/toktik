package service

import (
	token "github.com/Happy-Why/toktik-common/token"
	"github.com/Happy-Why/toktik-user/internal/global"
)

func CreateToken(userID int64) (string, string, error) {
	duration := global.PvSettings.Token.UserTokenExp
	data, err := token.NewTokenContent(token.UserToken, userID).Marshal()
	if err != nil {
		return "", "", err
	}
	tokenStr, _, err := global.Maker.CreateToken(data, duration)
	if err != nil {
		return "", "", err
	}
	return tokenStr, string(data), nil
}
