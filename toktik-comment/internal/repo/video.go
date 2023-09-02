package repo

import (
	"context"
	"toktik-comment/internal/model/auto"
)

type VideoRepo interface {
	CreateComment(c context.Context, commentInfo *auto.Comment) error
	DeleteComment(c context.Context, commentInfo *auto.Comment) error
	GetCommentAuthorIds(c context.Context, videoId int64) ([]int64, error)
	GetCommentList(c context.Context, videoId int64) ([]*auto.Comment, error)
}
