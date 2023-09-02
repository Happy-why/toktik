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
	IsVideoExist(c context.Context, videoId int64) (bool, error)
}
