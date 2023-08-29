package auto

import "strconv"

// Favorite 点赞表 /*
type Favorite struct {
	ID      uint `json:"id,string" gorm:"primarykey"`
	UserId  uint `json:"user_id,string" gorm:"not null;uniqueIndex:user_video"`
	VideoId uint `json:"video_id,string" gorm:"not null;uniqueIndex:user_video;index:video"`
}

func (*Favorite) TableName() string {
	return "favorite"
}

func CreateFavKey(userId uint) string {
	favStr := strconv.Itoa(int(userId))
	return "favorite::" + favStr
}
