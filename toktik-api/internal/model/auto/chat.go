package auto

type Message struct {
	BaseModel
	UserId   uint   `json:"from_user_id"`
	ToUserId uint   `json:"to_user_id"`
	Content  string `json:"content"`
}
