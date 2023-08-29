package auto

// Comment 评论表 /
type Comment struct {
	BaseModel
	VideoId uint   `json:"video_id,string" column:"video_id" gorm:"not null;index:comment_video"`
	UserId  uint   `json:"user_id,string" column:"user_id" gorm:"not null"`
	Content string `json:"content,string" column:"content"`
}

func (*Comment) TableName() string {
	return "comment"
}
