package auto

import (
	"encoding/json"
	"gorm.io/gorm"
	"strconv"
)

// Video 视频表 /*
type Video struct {
	gorm.Model
	UserId        uint   `json:"user_id" gorm:"not null;index"`
	Title         string `json:"title" gorm:"not null"`     // 视频标题
	PlayURL       string `json:"play_url" gorm:"not null"`  // 视频播放地址
	CoverURL      string `json:"cover_url" gorm:"not null"` // 视频封面地址
	FavoriteCount int64  `json:"favorite_count"`            // 视频的点赞总数
	CommentCount  int64  `json:"comment_count"`             // 视频的评论总数
	//CreatedAt time.Time `gorm:"index"`
}

func (*Video) TableName() string {
	return "video"
}

func NewVideoKey(videoId uint) string {
	videoStr := strconv.Itoa(int(videoId))
	return "video_info::" + videoStr
}

func NewMapVideoInfo(videoInfo *Video) map[string]interface{} {
	videoStr, _ := json.Marshal(videoInfo)
	videoMap := make(map[string]interface{})
	_ = json.Unmarshal(videoStr, &videoMap)
	return videoMap
}
