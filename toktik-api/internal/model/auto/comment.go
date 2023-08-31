package auto

// Comment 评论表 /
type Comment struct {
	BaseModel
	VideoId uint   `json:"video_id,string" gorm:"not null;uniqueIndex:comment_video"`
	UserId  uint   `json:"user_id,string" gorm:"not null"`
	Content string `json:"content,string" `
}

func (*Comment) TableName() string {
	return "comment"
}
