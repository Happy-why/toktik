package repo

import (
	"context"
	"github.com/Happy-Why/toktik-video/internal/model/auto"
)

type VideoRepo interface {
	CreateVideo(c context.Context, videoInfo *auto.Video) error
}
