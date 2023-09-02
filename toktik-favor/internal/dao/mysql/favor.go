package mysql

import (
	"context"
	"toktik-favor/internal/model/auto"
)

type FavorDao struct {
	conn *GormConn
}

func NewFavorDao() *FavorDao {
	return &FavorDao{
		conn: NewGormConn(),
	}
}

func (v *FavorDao) IsFavoriteVideo(c context.Context, userId, videoId int64) (bool, error) {
	var count int64
	err := v.conn.Session(c).Model(&auto.Favorite{}).Where("user_id = ? AND video_id = ?", userId, videoId).Count(&count).Error
	return count > 0, err
}

func (v *FavorDao) IsFavoriteRecordExist(c context.Context, userId, videoId int64) (bool, error) {
	var count int64
	err := v.conn.Session(c).Model(&auto.Favorite{}).
		Where("user_id = ? AND video_id = ?", userId, videoId).Count(&count).Error
	return count > 0, err
}

func (v *FavorDao) CreateFavoriteRecord(c context.Context, conn DbConn, favoriteInfo *auto.Favorite) error {
	v.conn = conn.(*GormConn)
	return v.conn.Tx(c).Create(favoriteInfo).Error
}

func (v *FavorDao) DeleteFavoriteRecord(c context.Context, conn DbConn, favoriteInfo *auto.Favorite) error {
	v.conn = conn.(*GormConn)
	return v.conn.Tx(c).Model(&auto.Favorite{}).
		Where("user_id = ? AND video_id = ?", favoriteInfo.UserId, favoriteInfo.VideoId).
		Unscoped().Delete(favoriteInfo).Error
}
