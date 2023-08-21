package mysql

import (
	"context"
	"github.com/Happy-Why/toktik-video/internal/model/auto"
)

type VideoDao struct {
	conn *GormConn
}

func NewVideoDao() *VideoDao {
	return &VideoDao{
		conn: NewGormConn(),
	}
}

func (v *VideoDao) CreateVideo(c context.Context, videoInfo *auto.Video) error {
	return v.conn.Session(c).Create(videoInfo).Error
}

func (v *VideoDao) GetVideosByTime(c context.Context, LatestTime int64) ([]*auto.Video, error) {
	videos := make([]*auto.Video, 0)
	session := v.conn.Session(c).Model(&auto.Video{}).Where("created_at < ?", LatestTime).Order("created_at DESC").Limit(30)
	err := session.Find(&videos).Error
	return videos, err
}

func (v *VideoDao) GetUserInfoByID(c context.Context, userId int64) (*auto.User, error) {
	userInfos := &auto.User{}
	err := v.conn.Session(c).Model(&auto.User{}).Where("id = ?", userId).First(userInfos).Error
	return userInfos, err
}

func (v *VideoDao) IsFollowTargetUser(c context.Context, userId, targetId int64) (bool, error) {
	var count int64
	err := v.conn.Session(c).Model(&auto.Relation{}).Where("user_id = ? AND target_id = ?", userId, targetId).Count(&count).Error
	return count > 0, err
}

func (v *VideoDao) IsFavoriteVideo(c context.Context, userId, videoId int64) (bool, error) {
	var count int64
	err := v.conn.Session(c).Model(&auto.Favorite{}).Where("user_id = ? AND video_id = ?", userId, videoId).Count(&count).Error
	return count > 0, err
}
