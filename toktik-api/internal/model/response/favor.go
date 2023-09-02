package response

type FavoriteActionResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type FavoriteListResponse struct {
	StatusCode int64    `json:"status_code"`
	StatusMsg  string   `json:"status_msg"`
	VideoList  []*Video `json:"video_list"`
}
