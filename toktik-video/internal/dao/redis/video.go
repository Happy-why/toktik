package redis

import (
	"context"
	"strconv"
	"strings"
	"time"
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
func (vr *VideoRClient) GetKeys(c context.Context, keyPattern string) ([]string, error) {
	return vr.rClient.GetKeys(c, keyPattern)
}
func (vr *VideoRClient) KeyExist(c context.Context, key string) (bool, error) {
	a, err := vr.rClient.KeyExist(c, key)
	if a == 0 {
		return true, err
	}
	return true, nil
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
	str := "2" + "+" + strconv.FormatInt(videoId, 10)
	_, err := vr.rClient.SDel(c, key, str)
	if err != nil {
		return err
	}
	videoStr := "1" + "+" + strconv.FormatInt(videoId, 10)
	_, err = vr.rClient.SAdd(c, key, videoStr)
	return err
}

func (vr *VideoRClient) DelFavorite(c context.Context, key string, videoId int64) error {
	str := "1" + "+" + strconv.FormatInt(videoId, 10)
	_, err := vr.rClient.SDel(c, key, str)
	if err != nil {
		return err
	}
	videoStr := "2" + "+" + strconv.FormatInt(videoId, 10)
	_, err = vr.rClient.SAdd(c, key, videoStr)
	return err
}

func (vr *VideoRClient) IsFavRecordExist(c context.Context, key string, videoId int64) (bool, error) {
	return vr.rClient.SIsExist(c, key, videoId)
}

func (vr *VideoRClient) HGetVideoInfo(c context.Context, key string) (*auto.Video, error) {
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
	if len(videoIdsStr) == 0 {
		return nil, nil
	}
	for _, v := range videoIdsStr {
		videoId, err = strconv.ParseInt(v, 10, 64)
		videoIds = append(videoIds, videoId)
	}
	return videoIds, err
}

func (vr *VideoRClient) PublishVideo(c context.Context, key string, time float64, videoId int64) error {
	_, err := vr.rClient.ZSet(c, key, []float64{time}, []interface{}{videoId})
	return err
}

func (vr *VideoRClient) ZGetVideoIds(c context.Context, key string, latestTime int64) (int64, []int64, error) {
	var videoId int64
	videoStr, err := vr.rClient.ZGetRevRangeWithScores(c, key, "-inf", strconv.FormatInt(latestTime, 10), 0, 30)
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

func (vr *VideoRClient) AddVideoCommentCount(c context.Context, key string) error {
	_, err := vr.rClient.IncrHMCount(c, key, auto.CommentCount, 1)
	return err
}

func (vr *VideoRClient) SubVideoCommentCount(c context.Context, key string) error {
	_, err := vr.rClient.IncrHMCount(c, key, auto.CommentCount, -1)
	return err
}

func (vr *VideoRClient) SDelUserVideo(c context.Context, key string) error {
	return vr.rClient.Del(c, key)
}

func (vr *VideoRClient) SGetUserVideoIds(c context.Context, key string) ([]int64, error) {
	str, err := vr.rClient.SGetAll(c, key)
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

func (vr *VideoRClient) SAddUserVideoIds(c context.Context, key string, videoId []int64) error {
	str := make([]string, len(videoId))
	for i, v := range videoId {
		str[i] = strconv.FormatInt(v, 10)
	}
	_, err := vr.rClient.SAdd(c, key, str)
	if err != nil {
		return err
	}
	_, err = vr.rClient.Expire(c, key, time.Hour*48)
	return err
}

func (vr *VideoRClient) ZSAddCommentInfo(c context.Context, key string, createTime float64, content string) error {
	_, err := vr.rClient.ZSet(c, key, []float64{createTime}, []interface{}{content})
	if err != nil {
		return err
	}
	_, err = vr.rClient.Expire(c, key, 48*time.Hour)
	return err
}

func (vr *VideoRClient) DelComment(c context.Context, key string) error {
	return vr.rClient.Del(c, key)
}

func (vr *VideoRClient) ZGetCommentList(c context.Context, key string) ([]int64, []*auto.Comment, error) {
	commentStrList, err := vr.rClient.ZGetRangeWithScores(c, key, 0, -1)
	if err != nil {
		return nil, nil, err
	}
	if len(commentStrList) == 0 {
		return nil, nil, err
	}
	userIds := make([]int64, len(commentStrList))
	//commentIds := make([]int64, len(commentStrList))
	//contents := make([]string, len(commentStrList))
	//creatTimes := make([]time.Time, len(commentStrList))
	commentList := make([]*auto.Comment, len(commentStrList))
	for i, v := range commentStrList {
		str := strings.SplitN(v.Member.(string), "+", 3)
		userIds[i], _ = strconv.ParseInt(str[1], 10, 64)
		commentId, _ := strconv.ParseInt(str[0], 10, 64)
		commentList = append(commentList, &auto.Comment{
			BaseModel: auto.BaseModel{ID: uint(commentId), CreatedAt: time.Unix(int64(v.Score), 0)},
			Content:   str[2],
		})
	}
	return userIds, commentList, nil
}
