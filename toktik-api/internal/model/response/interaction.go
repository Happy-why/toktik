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
	StatusCode int64   `json:"status_code"`
	StatusMsg  string  `json:"status_msg"`
	UserList   []*User `json:"user_list"`
}
