package request

type FollowActionRequest struct {
	Token      string `json:"token" form:"token" query:"token"`
	ToUserId   int64  `json:"to_user_id" form:"to_user_id" query:"to_user_id"`
	ActionType int32  `json:"action_type" form:"action_type" query:"action_type"`
	MyUserId   int64  `json:"my_user_id" form:"my_user_id" query:"my_user_id"`
}

type FollowListRequest struct {
	UserId   int64  `json:"user_id" form:"user_id" query:"user_id"`
	MyUserId int64  `json:"my_user_id" form:"my_user_id" query:"my_user_id"`
	Token    string `json:"token" form:"token" query:"token"`
}

type FansListRequest struct {
	UserId   int64  `json:"user_id" form:"user_id" query:"user_id"`
	MyUserId int64  `json:"my_user_id" form:"my_user_id" query:"my_user_id"`
	Token    string `json:"token" form:"token" query:"token"`
}

type FriendListRequest struct {
	UserId   int64  `json:"user_id" form:"user_id" query:"user_id"`
	MyUserId int64  `json:"my_user_id" form:"my_user_id" query:"my_user_id"`
	Token    string `json:"token" form:"token" query:"token"`
}
