package response

type ChatActionResponse struct {
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type MessageListResponse struct {
	StatusCode  int64              `json:"status_code"`
	StatusMsg   string             `json:"status_msg"`
	MessageList []*MessageResponse `json:"message_list"`
}

type MessageResponse struct {
	ID         uint   `gorm:"primarykey" json:"id,string"`
	CreateTime int64  `json:"create_time"`
	FromUserId uint   `json:"from_user_id"`
	ToUserId   uint   `json:"to_user_id"`
	Content    string `json:"content"`
}
