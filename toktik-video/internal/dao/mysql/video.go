package mysql

import (
	"context"
	"fmt"
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

func (v *VideoDao) IsFavoriteVideo(c context.Context, userId, videoId int64) (bool, error) {
	var count int64
	err := v.conn.Session(c).Model(&auto.Favorite{}).Where("user_id = ? AND video_id = ?", userId, videoId).Count(&count).Error
	return count > 0, err
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
func (v *VideoDao) IsFavoriteRecordExist(c context.Context, userId, videoId int64) (bool, error) {
	var count int64
	err := v.conn.Session(c).Model(&auto.Favorite{}).
		Where("user_id = ? AND video_id = ?", userId, videoId).Count(&count).Error
	return count > 0, err
}

func (v *VideoDao) CreateFavoriteRecord(c context.Context, conn DbConn, favoriteInfo *auto.Favorite) error {
	v.conn = conn.(*GormConn)
	return v.conn.Tx(c).Create(favoriteInfo).Error
}

func (v *VideoDao) AddVideoFavoriteCount(c context.Context, conn DbConn, videoId int64) error {
	v.conn = conn.(*GormConn)
	return v.conn.Tx(c).Model(&auto.Video{}).
		Where("id = ?", videoId).
		Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error
}
func (v *VideoDao) DeleteFavoriteRecord(c context.Context, conn DbConn, favoriteInfo *auto.Favorite) error {
	v.conn = conn.(*GormConn)
	return v.conn.Tx(c).Model(&auto.Favorite{}).
		Where("user_id = ? AND video_id = ?", favoriteInfo.UserId, favoriteInfo.VideoId).
		Unscoped().Delete(favoriteInfo).Error
}

func (v *VideoDao) SubVideoFavoriteCount(c context.Context, conn DbConn, videoId int64) error {
	v.conn = conn.(*GormConn)
	return v.conn.Tx(c).Model(&auto.Video{}).
		Where("id = ?", videoId).
		Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error
}

func (v *VideoDao) GetFavoriteVideosByUserId(c context.Context, userId int64) ([]*auto.Video, error) {
	videos := make([]*auto.Video, 0)
	session := v.conn.Session(c)
	sql := fmt.Sprintf("SELECT v.* from favorite f JOIN video v on f.video_id = v.id WHERE f.user_id = ?;")
	raw := session.Raw(sql, userId)
	err := raw.Scan(&videos).Error
	return videos, err
}

func (v *VideoDao) CreateComment(c context.Context, commentInfo *auto.Comment) error {
	return v.conn.Session(c).Create(commentInfo).Error
}

func (v *VideoDao) DeleteComment(c context.Context, commentInfo *auto.Comment) error {
	return v.conn.Session(c).Model(&auto.Comment{}).Where("id = ?", commentInfo.ID).Unscoped().Delete(commentInfo).Error
}

func (v *VideoDao) GetCommentAuthorIds(c context.Context, videoId int64) ([]int64, error) {
	userIds := make([]int64, 0)
	err := v.conn.Session(c).Model(&auto.Comment{}).
		Where("video_id = ?", videoId).
		Order("created_at desc").
		Pluck("user_id", &userIds).Error
	return userIds, err
}

func (v *VideoDao) GetCommentList(c context.Context, videoId int64) ([]*auto.Comment, error) {
	commentInfos := make([]*auto.Comment, 0)
	err := v.conn.Session(c).Model(&auto.Comment{}).
		Where("video_id = ?", videoId).
		Order("created_at desc").
		Find(&commentInfos).Error
	return commentInfos, err
}
