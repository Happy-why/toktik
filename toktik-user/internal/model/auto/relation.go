package auto

import "gorm.io/gorm"

type Relation struct {
	gorm.Model
	UserId   uint `json:"user_id" gorm:"not null"`   // 用户ID
	TargetId uint `json:"target_id" gorm:"not null"` // 目标ID
	//IsFriend int  `json:"is_friend" gorm:"not null"` // 如果需要保证 relation_id 唯一，可以使用该字段
}

func (*Relation) TableName() string {
	return "relation"
}

// 保证 relation_id 唯一，两个人中只要有一个人发起关注，就会创建一个 relation_id ，
// 在创建 relation_id，先查询 两人是否存在关系（比如 对方是否已经关注自己）
// 若对方没有关注自己，则直接创建 relation
// 对方已经关注自己，将 isFriend 值改为 1
