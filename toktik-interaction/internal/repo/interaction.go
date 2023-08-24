package repo

import (
	"context"
	"github.com/Happy-Why/toktik-interaction/internal/dao/mysql"
	"github.com/Happy-Why/toktik-interaction/internal/model/auto"
)

type InteractionRepo interface {
	IsRelationExist(c context.Context, myUserID, targetID int64) (uint, bool, error)
	FollowUserAction(c context.Context, conn mysql.DbConn, relationInfo *auto.Relation) error
	AddFollowCount(c context.Context, conn mysql.DbConn, userID uint) error
	AddFollowerCount(c context.Context, conn mysql.DbConn, userID uint) error
	CancelFollowUser(c context.Context, conn mysql.DbConn, relationInfo *auto.Relation) error
	SubFollowCount(c context.Context, conn mysql.DbConn, userID uint) error
	SubFollowerCount(c context.Context, conn mysql.DbConn, userID uint) error
	GetFollowIDs(c context.Context, userID uint) ([]int64, error)
	GetFansIDs(c context.Context, userID uint) ([]int64, error)
	GetFriendIDs(c context.Context, userID uint) ([]int64, error)
	IsFollowUser(c context.Context, myUserID, targetUserID int64) (bool, error)
}
