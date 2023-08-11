package mysql

import (
	"context"
	"github.com/Happy-Why/toktik-user/internal/model/auto"
)

type UserDao struct {
	conn *GormConn
}

func NewUserDao() *UserDao {
	return &UserDao{
		conn: NewGormConn(),
	}
}

func (u *UserDao) GetUserByUsername(c context.Context, username string) (bool, error) {
	var count int64
	err := u.conn.Session(c).Model(&auto.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}
