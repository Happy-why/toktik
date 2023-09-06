package auto

// Favorite 点赞表 /*
type Favorite struct {
	ID      uint `json:"id,string" gorm:"primarykey"`
	UserId  uint `json:"user_id,string" gorm:"uniqueIndex:user_video_id,not null;"`
	VideoId uint `json:"video_id,string" gorm:"uniqueIndex:user_video_id,not null;"`
}

func (*Favorite) TableName() string {
	return "favorite"
}
