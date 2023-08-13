package response

type RegisterResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	UserId     uint   `json:"user_id"`
	Token      string `json:"token"`
}

type LoginResponse struct {
	StatusCode int64  ` json:"status_code"`
	StatusMsg  string ` json:"status_msg"`
	UserId     int64  `json:"user_id"`
	Token      string ` json:"token"`
}

type UserIndexResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	User       *User  `json:"user"`
}

type User struct {
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
}
