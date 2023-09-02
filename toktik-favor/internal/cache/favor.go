package cache

import (
	"context"
	"strconv"
	"strings"
)

type FavorCache struct {
	rCache *RdbCache
}

func NewFavorCache() *FavorCache {
	return &FavorCache{
		rCache: GetRdbCache(),
	}
}

func (fc *FavorCache) GetKeys(c context.Context, keyPattern string) ([]string, error) {
	return fc.rCache.GetKeys(c, keyPattern)
}
func (fc *FavorCache) KeyExist(c context.Context, key string) (bool, error) {
	a, err := fc.rCache.KeyExist(c, key)
	if a == 0 {
		return false, err
	}
	return true, nil
}

func (fc *FavorCache) IsFavRecordExist(c context.Context, key string, videoId int64) (bool, error) {
	videoStr := "1+" + strconv.FormatInt(videoId, 10)
	return fc.rCache.SIsExist(c, key, videoStr)
}

func (fc *FavorCache) CreateFavorite(c context.Context, key string, videoId int64) error {
	str := "2" + "+" + strconv.FormatInt(videoId, 10)
	_, err := fc.rCache.SDel(c, key, str)
	if err != nil {
		return err
	}
	videoStr := "1" + "+" + strconv.FormatInt(videoId, 10)
	_, err = fc.rCache.SAdd(c, key, videoStr)
	return err
}

func (fc *FavorCache) DelFavorite(c context.Context, key string, videoId int64) error {
	str := "1" + "+" + strconv.FormatInt(videoId, 10)
	_, err := fc.rCache.SDel(c, key, str)
	if err != nil {
		return err
	}
	videoStr := "2" + "+" + strconv.FormatInt(videoId, 10)
	_, err = fc.rCache.SAdd(c, key, videoStr)
	return err
}

func (fc *FavorCache) GetFavoriteVideoIds(c context.Context, key string) ([]int64, error) {
	videoIds := make([]int64, 0)
	var videoId int64
	videoIdsStr, err := fc.rCache.SGetAll(c, key)
	if err != nil {
		return nil, err
	}
	if len(videoIdsStr) == 0 {
		return nil, nil
	}
	for _, v := range videoIdsStr {
		// 拆分 value
		str := strings.Split(v, "+")
		flag := str[0]
		videoId, _ = strconv.ParseInt(str[1], 10, 64)
		if flag == "2" {
			continue
		}
		videoIds = append(videoIds, videoId)
	}
	return videoIds, err
}
