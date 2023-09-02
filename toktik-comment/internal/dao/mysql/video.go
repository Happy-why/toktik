package mysql

import (
	"context"
	"toktik-comment/internal/model/auto"
)

type CommentDao struct {
	conn *GormConn
}

func NewCommentDao() *CommentDao {
	return &CommentDao{
		conn: NewGormConn(),
	}
}

func (v *CommentDao) CreateComment(c context.Context, commentInfo *auto.Comment) error {
	return v.conn.Session(c).Create(commentInfo).Error
}

func (v *CommentDao) DeleteComment(c context.Context, commentInfo *auto.Comment) error {
	return v.conn.Session(c).Model(&auto.Comment{}).Where("id = ?", commentInfo.ID).Unscoped().Delete(commentInfo).Error
}

func (v *CommentDao) GetCommentAuthorIds(c context.Context, videoId int64) ([]int64, error) {
	userIds := make([]int64, 0)
	err := v.conn.Session(c).Model(&auto.Comment{}).
		Where("video_id = ?", videoId).
		Order("created_at desc").
		Pluck("user_id", &userIds).Error
	return userIds, err
}

func (v *CommentDao) GetCommentList(c context.Context, videoId int64) ([]*auto.Comment, error) {
	commentInfos := make([]*auto.Comment, 0)
	err := v.conn.Session(c).Model(&auto.Comment{}).
		Where("video_id = ?", videoId).
		Order("created_at desc").
		Find(&commentInfos).Error
	return commentInfos, err
}
