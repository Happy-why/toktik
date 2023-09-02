package response

type CommentActionResponse struct {
	StatusCode int64    `json:"status_code"`
	StatusMsg  string   `json:"status_msg"`
	Comment    *Comment `json:"comment"`
}

type CommentListResponse struct {
	StatusCode  int64      `json:"status_code"`
	StatusMsg   string     `json:"status_msg"`
	CommentList []*Comment `json:"comment_list"`
}

type Comment struct {
	Id         int64  `json:"id"`          // 视频评论id
	User       *User  `json:"user"`        // 评论用户信息
	Content    string `json:"content"`     // 评论内容
	CreateDate string `json:"create_date"` // 评论发布日期，格式 mm-dd
}
