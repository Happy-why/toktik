package request

type ChatActionRequest struct {
	Token      string `json:"token" form:"token" query:"token"`
	UserId     int64  `json:"user_id" form:"user_id" query:"user_id"`
	ToUserId   int64  `json:"to_user_id" form:"to_user_id" query:"to_user_id"`
	ActionType int32  `json:"action_type" form:"action_type" query:"action_type"`
	Content    string `json:"content" form:"content" query:"content"`
}

func (r *ChatActionRequest) Verify() bool {
	return r.UserId == r.ToUserId
}

type MessageListRequest struct {
	UserId     int64  `json:"user_id" form:"user_id" query:"user_id"`
	Token      string `json:"token" form:"token" query:"token"`
	ToUserId   int64  `json:"to_user_id" form:"to_user_id" query:"to_user_id"`
	PreMsgTime int64  `json:"pre_msg_time" form:"pre_msg_time" query:"pre_msg_time"`
}
