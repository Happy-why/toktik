package mysql

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"testing"
	"time"
)

var Mdb *gorm.DB

func TestGetDB(t *testing.T) {
	InitMysql2()
	userID := 7093154536201650176
	UserIDs := make([]int64, 0)
	sql := "select target_id from relation where user_id = ?"
	raw := Mdb.Raw(sql, userID)
	err := raw.Scan(&UserIDs).Error
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Printf("userIDs:%#v\n", UserIDs[0])
	//err := Mdb.Model(&auto.Relation{}).Where("user_id = ?", userID).Scan(&UserIDs).Error
	//if err != nil {
	//	fmt.Println("err:", err)
	//}
	//fmt.Printf("userIDs:%#v\n", UserIDs[0])
}

func InitMysql2() {
	//配置MySQL连接参数
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "123456", "192.168.30.134", "3309", "toktik-interaction")
	fmt.Println(dsn)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢阙值
			Colorful:      true,        //禁用彩色
			LogLevel:      logger.Info,
		})
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false,
		},
	})
	if err != nil {
		zap.L().Error("init mysql db error:", zap.Error(err))
		panic("连接数据库失败, error=" + err.Error())
	}
	Mdb = DB
	//_ = DB.AutoMigrate(&User{})
}
