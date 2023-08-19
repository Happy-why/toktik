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
	GetFollowList(c context.Context, userID uint) ([]*auto.User, error)
	GetFansList(c context.Context, userID uint) ([]*auto.User, error)
	GetFriendList(c context.Context, userID uint) ([]*auto.User, error)
}
