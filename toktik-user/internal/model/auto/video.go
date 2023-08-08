package auto

import (
	"gorm.io/gorm"
)

// Video 视频表 /*
type Video struct {
	gorm.Model
	UserId        uint32 `json:"user_id" gorm:"not null;index"`
	Title         string `json:"title" gorm:"not null"`    // 视频标题
	PlayURL       string `json:"play_url" gorm:"not null"` // 视频播放地址
	CoverURL      string `json:"cover_url"`                // 视频封面地址
	FavoriteCount int64  `json:"favorite_count"`           // 视频的点赞总数
	CommentCount  int64  `json:"comment_count"`            // 视频的评论总数
	//CreatedAt time.Time `gorm:"index"`
}

func (*Video) TableName() string {
	return "video"
}
