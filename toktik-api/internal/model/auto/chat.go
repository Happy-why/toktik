package auto

type Message struct {
	BaseModel
	UserId   uint   `json:"from_user_id" gorm:"index:idx_message;not null"`
	ToUserId uint   `json:"to_user_id" gorm:"index:idx_message;not null"`
	Content  string `json:"content" gorm:"not null"`
}

func (*Message) TableName() string {
	return "messages"
}
