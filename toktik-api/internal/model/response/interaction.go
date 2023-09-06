package response

type FollowActionResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type FollowListResponse struct {
	StatusCode int64   `json:"status_code"`
	StatusMsg  string  `json:"status_msg"`
	UserList   []*User `json:"user_list"`
}

type FansListResponse struct {
	StatusCode int64   `json:"status_code"`
	StatusMsg  string  `json:"status_msg"`
	UserList   []*User `json:"user_list"`
}

type FriendListResponse struct {
	StatusCode int64         `json:"status_code"`
	StatusMsg  string        `json:"status_msg"`
	UserList   []*FriendUser `json:"user_list"`
}

type FriendUser struct {
	Id              int64   `json:"id"`
	Name            string  `json:"name"`
	FollowCount     *int64  `json:"follow_count"`
	FollowerCount   *int64  `json:"follower_count"`
	IsFollow        bool    `json:"is_follow"`
	Avatar          *string `json:"avatar"`
	BackgroundImage *string `json:"background_image"`
	Signature       *string `json:"signature"`
	TotalFavorited  *int64  `json:"total_favorited"`
	WorkCount       *int64  `json:"work_count"`
	FavoriteCount   *int64  `json:"favorite_count"`
	ChatMessage     string  `json:"message"`
	MsgType         int32   `json:"msgType"`
}
