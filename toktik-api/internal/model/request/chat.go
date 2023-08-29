package request

type ChatActionRequest struct {
	Token      string `json:"token" form:"token"`
	UserId     int64  `json:"user_id" form:"user_id"`
	ToUserId   int64  `json:"to_user_id" form:"to_user_id"`
	ActionType int32  `json:"action_type" form:"action_type"`
	Content    string `json:"content" form:"content"`
}

func (r *ChatActionRequest) Verify() bool {
	return r.UserId == r.ToUserId
}

type MessageListRequest struct {
	UserId   int64  `json:"user_id" form:"user_id"`
	Token    string `json:"token" form:"token"`
	ToUserId int64  `json:"to_user_id" form:"to_user_id"`
}
