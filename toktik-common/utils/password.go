package utils

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

/*
	原理:
	是一种加盐的加密方法，MD5加密时候，同一个密码经过hash的时候生成的是同一个hash值，在大数据的情况下，有些经过md5加密的方法将会被破解.
	使用BCrypt进行加密，同一个密码每次生成的hash值都是不相同的。
	每次加密的时候首先会生成一个随机数就是盐，之后将这个随机数与密码进行hash，得到一个hash值存到数据库。
	当用户在登陆的时候，输入的是明文的密码，
	从数据库中取出保存密码对其hash值进行分离，前面的22位就是加的盐，之后将随机数与前端输入的密码进行组合求hash值判断是否相同
*/

// HashPassword 计算bcrypt哈希字符串
func HashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("未能包装密码:%v", err)
	}
	return string(hashPassword), nil
}

// CheckPassword 检查输入的密码和哈希字符串是否匹配
func CheckPassword(password, hashPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}
