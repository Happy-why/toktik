package mysql

import (
	"context"
	"gorm.io/gorm"
	"toktik-video/internal/model/auto"
)

type VideoDao struct {
	conn *GormConn
}

func NewVideoDao() *VideoDao {
	return &VideoDao{
		conn: NewGormConn(),
	}
}

func (v *VideoDao) CreateVideo(c context.Context, conn DbConn, videoInfo *auto.Video) error {
	v.conn = conn.(*GormConn)
	return v.conn.Tx(c).Create(videoInfo).Error
}

func (v *VideoDao) GetVideosByTime(c context.Context, LatestTime int64) ([]*auto.Video, error) {
	videos := make([]*auto.Video, 0)
	session := v.conn.Session(c).Model(&auto.Video{}).Where("created_at < ?", LatestTime).Order("created_at DESC").Limit(30)
	err := session.Find(&videos).Error
	return videos, err
}

func (v *VideoDao) GetVideosByUserId(c context.Context, userId int64) ([]*auto.Video, error) {
	videos := make([]*auto.Video, 0)
	err := v.conn.Session(c).Model(&auto.Video{}).Where("user_id = ?", userId).Find(&videos).Error
	return videos, err
}

func (v *VideoDao) GetVideoIdsByUserId(c context.Context, userId int64) ([]int64, error) {
	videoIds := make([]int64, 0)
	err := v.conn.Session(c).Model(&auto.Video{}).Where("user_id = ?", userId).Pluck("id", &videoIds).Error
	return videoIds, err
}

func (v *VideoDao) GetVideoInfoByVideoId(c context.Context, videoId int64) (*auto.Video, error) {
	var videoInfo *auto.Video
	err := v.conn.Session(c).Model(&auto.Video{}).Where("id = ?", videoId).First(&videoInfo).Error
	return videoInfo, err
}

func (v *VideoDao) IsVideoExist(c context.Context, videoId int64) (bool, error) {
	var count int64
	err := v.conn.Session(c).Model(&auto.Video{}).Where("id = ?", videoId).Count(&count).Error
	return count > 0, err
}

func (v *VideoDao) AddVideoFavoriteCount(c context.Context, conn DbConn, videoId int64) error {
	v.conn = conn.(*GormConn)
	return v.conn.Tx(c).Model(&auto.Video{}).
		Where("id = ?", videoId).
		Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
}

func (v *VideoDao) SubVideoFavoriteCount(c context.Context, conn DbConn, videoId int64) error {
	v.conn = conn.(*GormConn)
	return v.conn.Tx(c).Model(&auto.Video{}).
		Where("id = ?", videoId).
		Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error
}
