package cache

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"toktik-comment/internal/model/auto"
)

type CommentCache struct {
	rCache *RdbCache
}

func NewCommentCache() *CommentCache {
	return &CommentCache{
		rCache: GetRdbCache(),
	}
}
func (vr *CommentCache) GetKeys(c context.Context, keyPattern string) ([]string, error) {
	return vr.rCache.GetKeys(c, keyPattern)
}
func (vr *CommentCache) KeyExist(c context.Context, key string) (bool, error) {
	a, err := vr.rCache.KeyExist(c, key)
	if a == 0 {
		return false, err
	}
	return true, nil
}

func (vr *CommentCache) ZSAddCommentInfo(c context.Context, key string, createTime float64, content string) error {
	_, err := vr.rCache.ZSet(c, key, []float64{createTime}, []interface{}{content})
	if err != nil {
		return err
	}
	_, err = vr.rCache.Expire(c, key, 48*time.Hour)
	return err
}

func (vr *CommentCache) DelComment(c context.Context, key string) error {
	return vr.rCache.Del(c, key)
}

func (vr *CommentCache) ZGetCommentList(c context.Context, key string) ([]int64, []*auto.Comment, error) {
	commentStrList, err := vr.rCache.ZGetRangeWithScores(c, key, 0, -1)
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
		commentList[i] = &auto.Comment{
			BaseModel: auto.BaseModel{ID: uint(commentId), CreatedAt: time.Unix(int64(v.Score), 0)},
			Content:   str[2],
		}
		fmt.Printf("commentInfo: commentId:%#v\n", commentId)
	}
	return userIds, commentList, nil
}
