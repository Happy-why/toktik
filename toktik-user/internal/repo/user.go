package repo

import (
	"context"
	"toktik-user/internal/dao/mysql"
	"toktik-user/internal/model/auto"
)

type UserRepo interface {
	GetUserByUsername(c context.Context, username string) (bool, error)
	SetUserInfo(c context.Context, conn mysql.DbConn, userInfo *auto.User) error
	SetUserCountInfo(c context.Context, conn mysql.DbConn, userCountInfo *auto.UserCount) error

	GetUserInfoByUsername(c context.Context, username string) (*auto.User, error)
	GetUserInfoByUserID(c context.Context, userID int64) (*auto.User, error)
	GetUserCountInfoByUserID(c context.Context, userID int64) (*auto.UserCount, error)
	AddFollowCount(c context.Context, conn mysql.DbConn, userID uint) error
	AddFollowerCount(c context.Context, conn mysql.DbConn, userID uint) error
	SubFollowCount(c context.Context, conn mysql.DbConn, userID uint) error
	SubFollowerCount(c context.Context, conn mysql.DbConn, userID uint) error
	GetUserList(c context.Context, userIDs []int64) ([]*auto.User, error)
	GetUserCountList(c context.Context, userIDs []int64) ([]*auto.UserCount, error)
	// IsFollowUser(c context.Context, myUserID, targetUserID int64) (bool, error)
}
