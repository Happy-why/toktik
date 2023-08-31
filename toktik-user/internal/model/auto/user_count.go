package auto

type UserCount struct {
	BaseModel
	UserId         uint  `json:"user_id,string" gorm:"uniqueIndex,not null"`
	User           User  `gorm:"ForeignKey:UserId"`
	FollowCount    int64 `json:"follow_count,string" `    // 关注总数
	FollowerCount  int64 `json:"follower_count,string" `  // 粉丝总数
	TotalFavorited int64 `json:"total_favorited,string" ` // 获赞数量
	WorkCount      int64 `json:"work_count,string" `      // 作品数
	FavoriteCount  int64 `json:"favorite_count,string" `  // 点赞总数
}

func (*UserCount) TableName() string {
	return "user_count"
}

const (
	FollowCount    = "follow_count"
	FollowerCount  = "follower_count"
	TotalFavorited = "total_favorited"
	WorkCount      = "work_count"
	FavoriteCount  = "favorite_count"
)
