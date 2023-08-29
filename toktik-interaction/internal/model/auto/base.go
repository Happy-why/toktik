package auto

import (
	"gorm.io/gorm"
	"strconv"
	"time"
)

type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id,string"`
	CreatedAt time.Time      `json:"created_at,string"`
	UpdatedAt time.Time      `json:"updated_at,string"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,string"`
}

func CreateUserKey(userId uint) string {
	userStr := strconv.Itoa(int(userId))
	return "user_info::" + userStr
}
