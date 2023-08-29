package repo

import (
	"context"
	"toktik-interaction/internal/dao/mysql"
	"toktik-interaction/internal/model/auto"
)

type InteractionRepo interface {
	IsRelationExist(c context.Context, myUserID, targetID int64) (uint, bool, error)
	FollowUserAction(c context.Context, conn mysql.DbConn, relationInfo *auto.Relation) error
	CancelFollowUser(c context.Context, conn mysql.DbConn, relationInfo *auto.Relation) error
	GetFollowIDs(c context.Context, userID uint) ([]int64, error)
	GetFansIDs(c context.Context, userID uint) ([]int64, error)
	GetFriendIDs(c context.Context, userID uint) ([]int64, error)
	IsFollowUser(c context.Context, myUserID, targetUserID int64) (bool, error)
}
