package repo

import (
	"context"
	"github.com/Happy-Why/toktik-video/internal/model/auto"
)

type VideoRepo interface {
	CreateVideo(c context.Context, videoInfo *auto.Video) error
	GetVideosByTime(c context.Context, LatestTime int64) ([]*auto.Video, error)
	GetUserInfoByID(c context.Context, userIds int64) (*auto.User, error)
	IsFollowTargetUser(c context.Context, userId, targetId int64) (bool, error)
	IsFavoriteVideo(c context.Context, userId, videoId int64) (bool, error)
}
