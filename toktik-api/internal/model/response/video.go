package response

type VideoFeedResponse struct {
	StatusCode int64    `json:"status_code"`
	StatusMsg  string   `json:"status_msg"`
	NextTime   int64    `json:"next_time"`
	VideoList  []*Video `json:"video_list"`
}

type VideoPublishResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type PublishListResponse struct {
	StatusCode int64    `json:"status_code"`
	StatusMsg  string   `json:"status_msg"`
	VideoList  []*Video `json:"video_list"`
}

type FavoriteActionResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type FavoriteListResponse struct {
	StatusCode int64    `json:"status_code"`
	StatusMsg  string   `json:"status_msg"`
	VideoList  []*Video `json:"video_list"`
}

type CommentActionResponse struct {
	StatusCode int64    `json:"status_code"`
	StatusMsg  string   `json:"status_msg"`
	Comment    *Comment `json:"comment"`
}

type CommentListResponse struct {
	StatusCode  int64      `json:"status_code"`
	StatusMsg   string     `json:"status_msg"`
	CommentList []*Comment `json:"comment"`
}

type Video struct {
	Id            int64  `json:"id"`             // 视频唯一标识
	Author        *User  `json:"author"`         // 视频作者信息
	PlayUrl       string `json:"play_url"`       // 视频播放地址
	CoverUrl      string `json:"cover_url"`      // 视频封面地址
	FavoriteCount int64  `json:"favorite_count"` // 视频的点赞总数
	CommentCount  int64  `json:"comment_count"`  // 视频的评论总数
	IsFavorite    bool   `json:"is_favorite"`    // true-已点赞，false-未点赞
	Title         string `json:"title"`          // 视频标题
}

type Comment struct {
	Id         int64  `json:"id"`          // 视频评论id
	User       *User  `json:"user"`        // 评论用户信息
	Content    string `json:"content"`     // 评论内容
	CreateDate string `json:"create_date"` // 评论发布日期，格式 mm-dd
}
