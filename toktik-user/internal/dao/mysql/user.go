package mysql

import (
	"context"
	"gorm.io/gorm"
	"toktik-user/internal/model/auto"
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

func (u *UserDao) SetUserInfo(c context.Context, conn DbConn, userInfo *auto.User) error {
	u.conn = conn.(*GormConn)
	return u.conn.Tx(c).Create(userInfo).Error
}

func (u *UserDao) SetUserCountInfo(c context.Context, conn DbConn, userCountInfo *auto.UserCount) error {
	u.conn = conn.(*GormConn)
	return u.conn.Tx(c).Create(userCountInfo).Error
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

func (u *UserDao) GetUserCountInfoByUserID(c context.Context, userID int64) (*auto.UserCount, error) {
	userCountInfo := &auto.UserCount{}
	err := u.conn.Session(c).Model(&auto.UserCount{}).Where("user_id = ?", userID).First(userCountInfo).Error
	return userCountInfo, err
}

//func (u *UserDao) IsFollowUser(c context.Context, myUserID, targetUserID int64) (bool, error) {
//	var count int64
//	err := u.conn.Session(c).Model(&auto.Relation{}).
//		Where("user_id = ? AND target_id = ?", myUserID, targetUserID).Count(&count).Error
//	return count > 0, err
//}

func (u *UserDao) AddFollowCount(c context.Context, conn DbConn, userID uint) error {
	u.conn = conn.(*GormConn)
	return u.conn.Tx(c).Model(&auto.UserCount{}).
		Where("user_id = ?", userID).
		Update("follow_count", gorm.Expr("follow_count + ?", 1)).Error
}

func (u *UserDao) AddFollowerCount(c context.Context, conn DbConn, userID uint) error {
	u.conn = conn.(*GormConn)
	return u.conn.Tx(c).Model(&auto.UserCount{}).
		Where("user_id = ?", userID).
		Update("follower_count", gorm.Expr("follower_count + ?", 1)).Error
}

func (u *UserDao) SubFollowCount(c context.Context, conn DbConn, userID uint) error {
	u.conn = conn.(*GormConn)
	return u.conn.Tx(c).Model(&auto.UserCount{}).
		Where("user_id = ?", userID).
		Update("follow_count", gorm.Expr("follow_count - ?", 1)).Error
}

func (u *UserDao) SubFollowerCount(c context.Context, conn DbConn, userID uint) error {
	u.conn = conn.(*GormConn)
	return u.conn.Tx(c).Model(&auto.UserCount{}).
		Where("user_id = ?", userID).
		Update("follower_count", gorm.Expr("follower_count - ?", 1)).Error
}

func (u *UserDao) GetUserList(c context.Context, userIDs []int64) ([]*auto.User, error) {
	var userInfos []*auto.User
	err := u.conn.Session(c).Model(&auto.User{}).Where("id IN ?", userIDs).Find(&userInfos).Error
	return userInfos, err
}

func (u *UserDao) GetUserCountList(c context.Context, userIDs []int64) ([]*auto.UserCount, error) {
	var userCountInfos []*auto.UserCount
	err := u.conn.Session(c).Model(&auto.UserCount{}).Where("id IN ?", userIDs).Find(&userCountInfos).Error
	return userCountInfos, err
}
