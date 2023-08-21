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
