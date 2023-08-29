package auto

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Video 视频表 /*
type Video struct {
	BaseModel
	UserId        uint   `json:"user_id,string" gorm:"not null;index"`
	Title         string `json:"title" gorm:"not null"`     // 视频标题
	PlayURL       string `json:"play_url" gorm:"not null"`  // 视频播放地址
	CoverURL      string `json:"cover_url" gorm:"not null"` // 视频封面地址
	FavoriteCount int64  `json:"favorite_count,string"`     // 视频的点赞总数
	CommentCount  int64  `json:"comment_count,string"`      // 视频的评论总数
	//CreatedAt time.Time `gorm:"index"`
}

func (*Video) TableName() string {
	return "video"
}

const (
	AuthorId      = "user_id"
	FavoriteCount = "favorite_count"
	CommentCount  = "comment_count"
)

func CreateVideoKey(videoId uint) string {
	videoStr := strconv.Itoa(int(videoId))
	return "video_info::" + videoStr
}

func CreateMapVideoInfo(videoInfo *Video) map[string]interface{} {
	videoStr, _ := json.Marshal(videoInfo)
	videoMap := make(map[string]interface{})
	_ = json.Unmarshal(videoStr, &videoMap)
	delete(videoMap, "created_at")
	delete(videoMap, "updated_at")
	delete(videoMap, "deleted_at")
	fmt.Println("userMap:", videoMap)
	return videoMap
}

func CreateVideoInfo(videoMap map[string]string) (*Video, error) {
	videoStr, _ := json.Marshal(videoMap)
	videoInfo := new(Video)
	err := json.Unmarshal(videoStr, videoInfo)
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("videoInfo:", videoInfo)
	return videoInfo, err
}
