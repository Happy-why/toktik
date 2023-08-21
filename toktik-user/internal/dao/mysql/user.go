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

func (u *UserDao) UserRegister(c context.Context, userInfo *auto.User) error {
	return u.conn.Session(c).Create(userInfo).Error
}

func (u *UserDao) GetUserInfoByUsername(c context.Context, username string) (*auto.User, error) {
	userInfo := &auto.User{}
	err := u.conn.Session(c).Model(&auto.User{}).Where("username = ?", username).First(userInfo).Error
	return userInfo, err
}

func (u *UserDao) GetUserInfoByUserID(c context.Context, userID int64) (*auto.User, error) {
	userInfo := &auto.User{}
	err := u.conn.Session(c).Model(&auto.User{}).Where("id = ?", userID).First(userInfo).Error
	return userInfo, err
}

func (u *UserDao) IsFollowUser(c context.Context, myUserID, targetUserID int64) (bool, error) {
	var count int64
	err := u.conn.Session(c).Model(&auto.Relation{}).
		Where("user_id = ? AND target_id = ?", myUserID, targetUserID).Count(&count).Error
	return count > 0, err
}
