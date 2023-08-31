package auto

type UserCount struct {
	BaseModel
	UserId         uint  `json:"user_id,string" gorm:"uniqueIndex,not null"`
	User           User  `gorm:"ForeignKey:UserId"`
	FollowCount    int64 `json:"follow_count,string" gorm:"not null"`    // 关注总数
	FollowerCount  int64 `json:"follower_count,string" gorm:"not null"`  // 粉丝总数
	TotalFavorited int64 `json:"total_favorited,string" gorm:"not null"` // 获赞数量
	WorkCount      int64 `json:"work_count,string" gorm:"not null"`      // 作品数
	FavoriteCount  int64 `json:"favorite_count,string" gorm:"not null"`  // 点赞总数
}

func (*UserCount) TableName() string {
	return "user_count"
}
