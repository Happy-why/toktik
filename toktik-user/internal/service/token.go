package service

import (
	token "toktik-common/token"
	"toktik-user/internal/global"
)

func CreateToken(userID int64) (string, string, error) {
	duration := global.Settings.Token.UserTokenExp
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
