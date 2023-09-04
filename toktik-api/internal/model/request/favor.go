package request

type FavoriteActionRequest struct {
	Token      string `json:"token" form:"token" query:"token"`
	VideoId    int64  `json:"video_id" form:"video_id" query:"video_id"`
	ActionType int32  `json:"action_type" form:"action_type" query:"action_type"`
	UserId     int64  `json:"user_id" form:"user_id" query:"user_id"`
}

type FavoriteListRequest struct {
	UserId   int64  `json:"user_id" form:"user_id" query:"user_id"`
	MyUserId int64  `json:"my_user_id" form:"my_user_id" query:"my_user_id"`
	Token    string `json:"token" form:"token" query:"token"`
}
