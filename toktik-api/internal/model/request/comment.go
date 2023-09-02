package request

type CommentActionRequest struct {
	UserId      int64  `json:"user_id" form:"user_id"`
	Token       string `json:"token" form:"token"`
	VideoId     int64  `json:"video_id" form:"video_id"`
	ActionType  int32  `json:"action_type" form:"action_type"`
	CommentText string `json:"comment_text" form:"comment_text"`
	CommentId   int64  `json:"comment_id" form:"comment_id"`
}

type CommentListRequest struct {
	Token   string `json:"token" form:"token"`
	VideoId int64  `json:"video_id,omitempty" form:"video_id"`
	UserId  int64  `json:"user_id" form:"user_id"`
}
