package repo

import (
	"context"
	"toktik-favor/internal/dao/mysql"
	"toktik-favor/internal/model/auto"
)

type FavorRepo interface {
	IsFavoriteVideo(c context.Context, userId, videoId int64) (bool, error)
	IsFavoriteRecordExist(c context.Context, userId, videoId int64) (bool, error)
	CreateFavoriteRecord(c context.Context, conn mysql.DbConn, favoriteInfo *auto.Favorite) error
	DeleteFavoriteRecord(c context.Context, conn mysql.DbConn, favoriteInfo *auto.Favorite) error
}
