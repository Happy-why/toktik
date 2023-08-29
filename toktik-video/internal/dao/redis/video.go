package redis

import (
	"context"
	"strconv"
	"toktik-video/internal/global"
	"toktik-video/internal/model/auto"
)

type VideoRClient struct {
	rClient *RdbCache
}

func NewVideoRClient() *VideoRClient {
	return &VideoRClient{
		rClient: GetRdbCache(),
	}
}

func (vr *VideoRClient) HSetVideoInfo(c context.Context, key string, value map[string]interface{}) error {
	err := vr.rClient.HSet(c, key, value)
	if err != nil {
		return err
	}
	_, err = vr.rClient.Expire(c, key, global.Settings.Rules.VideoInfoCacheExpire)
	return err
}

func (vr *VideoRClient) AddVideoFavoriteCount(c context.Context, key string) error {
	_, err := vr.rClient.IncrHMCount(c, key, auto.FavoriteCount, 1)
	return err
}

func (vr *VideoRClient) SubVideoFavoriteCount(c context.Context, key string) error {
	_, err := vr.rClient.IncrHMCount(c, key, auto.FavoriteCount, -1)
	return err
}

func (vr *VideoRClient) CreateFavorite(c context.Context, key string, videoId int64) error {
	_, err := vr.rClient.SAdd(c, key, videoId)
	return err
}

func (vr *VideoRClient) DelFavorite(c context.Context, key string, videoId int64) error {
	_, err := vr.rClient.SDel(c, key, videoId)
	return err
}

func (vr *VideoRClient) IsFavRecordExist(c context.Context, key string, videoId int64) (bool, error) {
	return vr.rClient.SIsExist(c, key, videoId)
}

func (vr *VideoRClient) GetVideoInfo(c context.Context, key string) (*auto.Video, error) {
	videoMapInfo, err := vr.rClient.HGetAll(c, key)
	if len(videoMapInfo) == 0 {
		return nil, err
	}
	videoInfo, _ := auto.CreateVideoInfo(videoMapInfo)
	return videoInfo, err
}

func (vr *VideoRClient) GetFavoriteVideoIds(c context.Context, key string) ([]int64, error) {
	videoIds := make([]int64, 0)
	var videoId int64
	videoIdsStr, err := vr.rClient.SGetAll(c, key)
	if err != nil {
		return nil, err
	}
	for _, v := range videoIdsStr {
		videoId, err = strconv.ParseInt(v, 10, 64)
		videoIds = append(videoIds, videoId)
	}
	return videoIds, err
}

func (vr *VideoRClient) PublishVideo(c context.Context, key string, time float64, url string) error {
	_, err := vr.rClient.ZSet(c, key, []float64{time}, []interface{}{url})
	return err
}

func (vr *VideoRClient) ZGetVideoIds(c context.Context, key string, latestTime int64) (int64, []int64, error) {
	var videoId int64
	videoStr, err := vr.rClient.ZGetRevRangeWithScores(c, key, "-inf", strconv.FormatInt(latestTime, 10), 0, 30)
	if err != nil {
		return 0, nil, err
	}
	videoIds := make([]int64, len(videoStr))
	for i, v := range videoStr {
		videoId, _ = strconv.ParseInt(v.Member.(string), 10, 64)
		videoIds[i] = videoId
	}
	nextTime := int64(videoStr[len(videoStr)-1].Score)
	return nextTime, videoIds, err
}

func (vr *VideoRClient) AddVideoCommentCount(c context.Context, key string) error {
	_, err := vr.rClient.IncrHMCount(c, key, auto.CommentCount, 1)
	return err
}

func (vr *VideoRClient) SubVideoCommentCount(c context.Context, key string) error {
	_, err := vr.rClient.IncrHMCount(c, key, auto.CommentCount, -1)
	return err
}
