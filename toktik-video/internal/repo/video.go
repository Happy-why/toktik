package repo

import (
	"context"
	"toktik-video/internal/dao/mysql"
	"toktik-video/internal/model/auto"
)

type VideoRepo interface {
	CreateVideo(c context.Context, conn mysql.DbConn, videoInfo *auto.Video) error
	GetVideoInfoByVideoId(c context.Context, videoId int64) (*auto.Video, error)
	GetVideosByTime(c context.Context, LatestTime int64) ([]*auto.Video, error)
	GetVideosByUserId(c context.Context, userId int64) ([]*auto.Video, error)
	GetVideoIdsByUserId(c context.Context, userId int64) ([]int64, error)
	IsFavoriteVideo(c context.Context, userId, videoId int64) (bool, error)
	IsVideoExist(c context.Context, videoId int64) (bool, error)
	IsFavoriteRecordExist(c context.Context, userId, videoId int64) (bool, error)
	CreateFavoriteRecord(c context.Context, conn mysql.DbConn, favoriteInfo *auto.Favorite) error
	DeleteFavoriteRecord(c context.Context, conn mysql.DbConn, favoriteInfo *auto.Favorite) error
	GetFavoriteVideosByUserId(c context.Context, userId int64) ([]*auto.Video, error)
	CreateComment(c context.Context, commentInfo *auto.Comment) error
	DeleteComment(c context.Context, commentInfo *auto.Comment) error
	GetCommentAuthorIds(c context.Context, videoId int64) ([]int64, error)
	GetCommentList(c context.Context, videoId int64) ([]*auto.Comment, error)
}
