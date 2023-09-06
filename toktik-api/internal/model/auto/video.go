package auto

// Video 视频表 /*
type Video struct {
	BaseModel
	UserId        uint   `json:"user_id,string" gorm:"index:user_id,not null"`
	Title         string `json:"title" gorm:"type:varchar(255);not null"`     // 视频标题
	PlayURL       string `json:"play_url" gorm:"type:varchar(255);not null"`  // 视频播放地址
	CoverURL      string `json:"cover_url" gorm:"type:varchar(255);not null"` // 视频封面地址
	FavoriteCount int64  `json:"favorite_count,string" `                      // 视频的点赞总数
	CommentCount  int64  `json:"comment_count,string"`                        // 视频的评论总数
}

func (*Video) TableName() string {
	return "video"
}
