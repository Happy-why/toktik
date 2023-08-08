package auto

import "gorm.io/gorm"

// Favorite 点赞表 /*
type Favorite struct {
	gorm.Model
	UserId  uint32 `json:"user_id" gorm:"not null;uniqueIndex:user_video"`
	VideoId uint32 `json:"video_id" gorm:"not null;uniqueIndex:user_video;index:video"`
}

func (*Favorite) TableName() string {
	return "favorite"
}
