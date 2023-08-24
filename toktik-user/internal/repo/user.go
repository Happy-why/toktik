package repo

import (
	"context"
	"github.com/Happy-Why/toktik-user/internal/dao/mysql"
	"github.com/Happy-Why/toktik-user/internal/model/auto"
)

type UserRepo interface {
	GetUserByUsername(c context.Context, username string) (bool, error)
	UserRegister(c context.Context, userInfo *auto.User) error
	GetUserInfoByUsername(c context.Context, username string) (*auto.User, error)
	GetUserInfoByUserID(c context.Context, userID int64) (*auto.User, error)
	AddFollowCount(c context.Context, conn mysql.DbConn, userID uint) error
	AddFollowerCount(c context.Context, conn mysql.DbConn, userID uint) error
	SubFollowCount(c context.Context, conn mysql.DbConn, userID uint) error
	SubFollowerCount(c context.Context, conn mysql.DbConn, userID uint) error
	GetUserList(c context.Context, userIDs []int64) ([]*auto.User, error)
	// IsFollowUser(c context.Context, myUserID, targetUserID int64) (bool, error)
}
