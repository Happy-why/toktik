package cache

import (
	"context"
	"strconv"
	"time"
	"toktik-video/internal/global"
	"toktik-video/internal/model/auto"
)

type VideoCache struct {
	rCache *RdbCache
}

func NewVideoCache() *VideoCache {
	return &VideoCache{
		rCache: GetRdbCache(),
	}
}
func (vr *VideoCache) GetKeys(c context.Context, keyPattern string) ([]string, error) {
	return vr.rCache.GetKeys(c, keyPattern)
}
func (vr *VideoCache) KeyExist(c context.Context, key string) (bool, error) {
	a, err := vr.rCache.KeyExist(c, key)
	if a == 0 {
		return false, err
	}
	return true, nil
}
func (vr *VideoCache) HSetVideoInfo(c context.Context, key string, value map[string]interface{}) error {
	err := vr.rCache.HSet(c, key, value)
	if err != nil {
		return err
	}
	_, err = vr.rCache.Expire(c, key, global.Settings.Rules.VideoInfoCacheExpire)
	return err
}

func (vr *VideoCache) AddVideoFavoriteCount(c context.Context, key string) error {
	_, err := vr.rCache.IncrHMCount(c, key, auto.FavoriteCount, 1)
	return err
}

func (vr *VideoCache) SubVideoFavoriteCount(c context.Context, key string) error {
	_, err := vr.rCache.IncrHMCount(c, key, auto.FavoriteCount, -1)
	return err
}

func (vr *VideoCache) HGetVideoInfo(c context.Context, key string) (*auto.Video, error) {
	videoMapInfo, err := vr.rCache.HGetAll(c, key)
	if len(videoMapInfo) == 0 {
		return nil, err
	}
	videoInfo, _ := auto.CreateVideoInfo(videoMapInfo)
	return videoInfo, err
}

func (vr *VideoCache) PublishVideo(c context.Context, key string, time float64, videoId int64) error {
	_, err := vr.rCache.ZSet(c, key, []float64{time}, []interface{}{videoId})
	return err
}

func (vr *VideoCache) ZGetVideoIds(c context.Context, key string, latestTime int64) (int64, []int64, error) {
	var videoId int64
	videoStr, err := vr.rCache.ZGetRevRangeByScoreWithScores(c, key, "-inf", strconv.FormatInt(latestTime, 10), 0, 30)
	if err != nil {
		return 0, nil, err
	}
	if len(videoStr) == 0 {
		return 0, nil, err
	}
	videoIds := make([]int64, len(videoStr))
	for i, v := range videoStr {
		videoId, _ = strconv.ParseInt(v.Member.(string), 10, 64)
		videoIds[i] = videoId
	}
	nextTime := int64(videoStr[0].Score)
	return nextTime, videoIds, err
}

func (vr *VideoCache) AddVideoCommentCount(c context.Context, key string) error {
	_, err := vr.rCache.IncrHMCount(c, key, auto.CommentCount, 1)
	return err
}

func (vr *VideoCache) SubVideoCommentCount(c context.Context, key string) error {
	_, err := vr.rCache.IncrHMCount(c, key, auto.CommentCount, -1)
	return err
}

func (vr *VideoCache) SDelUserVideo(c context.Context, key string) error {
	return vr.rCache.Del(c, key)
}

func (vr *VideoCache) SGetUserVideoIds(c context.Context, key string) ([]int64, error) {
	str, err := vr.rCache.SGetAll(c, key)
	if err != nil {
		return nil, err
	}
	if len(str) == 0 {
		return nil, nil
	}
	Ids := make([]int64, len(str))
	for i, v := range str {
		Ids[i], _ = strconv.ParseInt(v, 10, 64)
	}
	return Ids, err
}

func (vr *VideoCache) SAddUserVideoIds(c context.Context, key string, videoId []int64) error {
	str := make([]string, len(videoId))
	for i, v := range videoId {
		str[i] = strconv.FormatInt(v, 10)
	}
	_, err := vr.rCache.SAdd(c, key, str)
	if err != nil {
		return err
	}
	_, err = vr.rCache.Expire(c, key, time.Hour*48)
	return err
}
