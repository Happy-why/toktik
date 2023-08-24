package auto

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type User struct {
	BaseModel
	Username        string `json:"username" gorm:"type:varchar(200);;not null;uniqueIndex:user_name"`
	Password        string `json:"password" gorm:"not null;"`
	Avatar          string `json:"avatar"`                                // 用户头像
	FollowCount     int64  `json:"follow_count,string" gorm:"not null"`   // 关注总数
	FollowerCount   int64  `json:"follower_count,string" gorm:"not null"` // 粉丝总数
	BackgroundImage string `json:"background_image"`                      //背景图片
	IsFollow        bool   `json:"is_follow,string"`                      // true-已关注，false-未关注 在数据库中这个字段没用
	Signature       string `json:"signature"`                             // 个人简介
	TotalFavorited  int64  `json:"total_favorited,string"`                // 获赞数量
	WorkCount       int64  `json:"work_count,string"`                     // 作品数
	FavoriteCount   int64  `json:"favorite_count,string"`                 // 点赞总数
}

func (*User) TableName() string {
	return "user"
}

const (
	FollowCount    = "follow_count"
	FollowerCount  = "follower_count"
	TotalFavorited = "total_favorited"
	WorkCount      = "work_count"
	FavoriteCount  = "favorite_count"
)

func CreateUserKey(videoId uint) string {
	videoStr := strconv.Itoa(int(videoId))
	return "user_info::" + videoStr
}

func CreateMapUserInfo(userInfo *User) map[string]interface{} {
	userStr, _ := json.Marshal(userInfo)
	userMap := make(map[string]interface{})
	_ = json.Unmarshal(userStr, &userMap)
	delete(userMap, "created_at")
	delete(userMap, "updated_at")
	delete(userMap, "deleted_at")
	fmt.Println("userMap:", userMap)
	return userMap
}

func CreateUserInfo(userMap map[string]string) (*User, error) {
	userStr, _ := json.Marshal(userMap)
	userInfo := new(User)
	err := json.Unmarshal(userStr, userInfo)
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("userInfo:", userInfo)
	return userInfo, err
}
