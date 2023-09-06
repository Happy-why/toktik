package auto

type User struct {
	BaseModel
	Username        string `json:"username" gorm:"uniqueIndex:username;size:40;not null"` // 设置唯一索引，判断用户名是否重复
	Password        string `json:"password" gorm:"type:varchar(50);not null"`
	Avatar          string `json:"avatar" gorm:"type:varchar(255);not null"`           // 用户头像
	BackgroundImage string `json:"background_image" gorm:"type:varchar(255);not null"` //背景图片
	IsFollow        bool   `json:"is_follow,string" gorm:"not null"`                   // true-已关注，false-未关注 在数据库中这个字段没用
	Signature       string `json:"signature" gorm:"type:varchar(255);"`                // 个人简介
}

func (*User) TableName() string {
	return "user"
}
