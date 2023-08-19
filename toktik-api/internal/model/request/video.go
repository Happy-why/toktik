package request

import "time"

type VideoFeedRequest struct {
	LatestTime time.Time `json:"latest_time"`
	Token      string    `json:"token"`
}

type VideoPublishRequest struct {
	Data  string `json:"data" form:"data"`
	Token string `json:"token" form:"token"`
	Title string `json:"title" form:"title"`
}

type PublishListRequest struct {
	Token  string `json:"token"`
	UserId int64  `json:"user_id"`
}

type FavoriteActionRequest struct {
	Token      string `json:"token"`
	VideoId    int64  `json:"video_id"`
	ActionType int32  `json:"action_type"`
}

type FavoriteListRequest struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type CommentActionRequest struct {
	Token       string `json:"token"`
	VideoId     int64  `json:"video_id"`
	ActionType  int32  `json:"action_type"`
	CommentText string `json:"comment_text"`
	CommentId   int64  `json:"comment_id"`
}

type CommentListRequest struct {
	Token   string `json:"token"`
	VideoId int64  `json:"video_id,omitempty"`
}
