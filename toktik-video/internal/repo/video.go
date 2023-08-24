package repo

import (
	"context"
	"github.com/Happy-Why/toktik-video/internal/dao/mysql"
	"github.com/Happy-Why/toktik-video/internal/model/auto"
)

type VideoRepo interface {
	CreateVideo(c context.Context, conn mysql.DbConn, videoInfo *auto.Video) error
	GetVideosByTime(c context.Context, LatestTime int64) ([]*auto.Video, error)
	GetUserInfoByID(c context.Context, userIds int64) (*auto.User, error)
	IsFollowTargetUser(c context.Context, userId, targetId int64) (bool, error)
	IsFavoriteVideo(c context.Context, userId, videoId int64) (bool, error)
	AddUserWorkCount(c context.Context, conn mysql.DbConn, userId int64) error
	GetVideosByUserId(c context.Context, userId int64) ([]*auto.Video, error)
	IsVideoExist(c context.Context, videoId int64) (bool, error)
	IsFavoriteRecordExist(c context.Context, userId, videoId int64) (bool, error)
	CreateFavoriteRecord(c context.Context, conn mysql.DbConn, favoriteInfo *auto.Favorite) error
	AddUserFavoriteCount(c context.Context, conn mysql.DbConn, userId int64) error
	AddVideoFavoriteCount(c context.Context, conn mysql.DbConn, videoId int64) error
	DeleteFavoriteRecord(c context.Context, conn mysql.DbConn, favoriteInfo *auto.Favorite) error
	SubUserFavoriteCount(c context.Context, conn mysql.DbConn, userId int64) error
	SubVideoFavoriteCount(c context.Context, conn mysql.DbConn, videoId int64) error
	GetFavoriteVideosByUserId(c context.Context, userId int64) ([]*auto.Video, error)
}
