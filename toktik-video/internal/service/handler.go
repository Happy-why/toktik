package service

import (
	"context"
	"github.com/Happy-Why/toktik-rpc/kitex_gen/video"
	"github.com/Happy-Why/toktik-video/internal/dao"
	"github.com/Happy-Why/toktik-video/internal/dao/mysql"
	"github.com/Happy-Why/toktik-video/internal/repo"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct {
	cache       repo.Cache
	videoRepo   repo.VideoRepo
	respRepo    repo.HandlerResp
	transaction mysql.Transaction
}

func NewVideoService() *VideoServiceImpl {
	return &VideoServiceImpl{
		cache:       dao.Group.Rdb,
		videoRepo:   mysql.NewVideoDao(),
		respRepo:    repo.NewHandlerResps(),
		transaction: mysql.NewTransaction(),
	}
}

func (vs *VideoServiceImpl) VideoFeed(ctx context.Context, req *video.VideoFeedRequest) (resp *video.VideoFeedResponse, err error) {

	return nil, nil
}

func (vs *VideoServiceImpl) VideoPublish(ctx context.Context, req *video.VideoPublishRequest) (resp *video.VideoPublishResponse, err error) {
	return nil, nil
}

func (vs *VideoServiceImpl) PublishList(ctx context.Context, req *video.PublishListRequest) (resp *video.PublishListResponse, err error) {
	return nil, nil
}

func (vs *VideoServiceImpl) FavoriteList(ctx context.Context, req *video.FavoriteListRequest) (resp *video.FavoriteListResponse, err error) {
	return nil, nil
}

func (vs *VideoServiceImpl) FavoriteAction(ctx context.Context, req *video.FavoriteActionRequest) (resp *video.FavoriteActionResponse, err error) {
	return nil, nil
}

func (vs *VideoServiceImpl) CommentAction(ctx context.Context, req *video.CommentActionRequest) (resp *video.CommentActionResponse, err error) {
	return nil, nil
}

func (vs *VideoServiceImpl) CommentList(ctx context.Context, req *video.CommentListRequest) (resp *video.CommentListResponse, err error) {
	return nil, nil
}
