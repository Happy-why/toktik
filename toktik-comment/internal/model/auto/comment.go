package auto

import "strconv"

// Comment 评论表 /
type Comment struct {
	BaseModel
	VideoId uint   `json:"video_id,string" gorm:"not null;index:comment_video"`
	UserId  uint   `json:"user_id,string" gorm:"not null"`
	Content string `json:"content,string" gorm:"not null"`
}

func (*Comment) TableName() string {
	return "comment"
}

// Key Value: comment_id+content

func CreateCommentKey(videoId int64) string {
	return "video_comment::" + strconv.FormatInt(videoId, 10)
}

func CreateCommentValue(commentId, userId int64, content string) string {
	return strconv.FormatInt(commentId, 10) + "+" + strconv.FormatInt(userId, 10) + "+" + content
}
