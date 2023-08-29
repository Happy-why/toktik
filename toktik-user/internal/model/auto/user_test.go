package auto

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"

	"testing"
)

var rdb *redis.Client
var Mdb *gorm.DB

func TestCreateUserInfo2(t *testing.T) {

}

func TestCreateUserInfo(t *testing.T) {
	rdb = InitRedis()
	InitMysql()
	c := context.Background()
	key := "video_publish"
	a := "1693264714"
	result, err := rdb.ZRevRangeByScore(c, key, &redis.ZRangeBy{Min: "-inf", Max: a, Offset: 0, Count: 30}).Result()
	if err != nil {
		fmt.Println("rdb.LPush err:", err)
	}
	fmt.Println("result:", result)
	//key := "abc$qwe$zxc$666"
	//a := strings.SplitN(key, "$", 3)
	//fmt.Println(a)

	//result, err := rdb.LPopCount(c, key, 3).Result()
	//// 返回list中还剩多少数据
	//if err != nil {
	//	fmt.Println("rdb.LPush err:", err)
	//}
	//fmt.Println("result:", result)
	//f := []float64{1.0, 2.0, 50.0, 100.0}
	//b := []interface{}{"qwe", "asd", 999, 666}
	//z := make([]*redis.Z, 0)
	//for i := 0; i < len(f); i++ {
	//	z = append(z, &redis.Z{Score: f[i], Member: b[i]})
	//}
	//result, err := rdb.ZAdd(c, key, z...).Result()
	//if err != nil {
	//	fmt.Println("rdb.HSet err:", err)
	//}
	//fmt.Println("result:", result)
	//result, err := rdb.ZAdd(c, key, lang...).Result()
	//min := "-inf"
	//max := strconv.FormatInt(time.Now().Unix(), 10)
	//offset := int64(0)
	//count := int64(30)
	//result, err := rdb.ZRevRangeByScore(c, key, &redis.ZRangeBy{Min: min, Max: max, Offset: offset, Count: count}).Result()
	//if err != nil {
	//	fmt.Println("rdb.HSet err:", err)
	//}
	//fmt.Println("result:", result)
	//fmt.Println("result:", result)
	//key := "video_info::7093525438898634752"
	//result, err := rdb.HGet(c, key, "user_id").Result()
	//if err == redis.Nil {
	//	fmt.Println("err:", err)
	//}
	//if err != nil {
	//	fmt.Println("rdb.HSet err:", err)
	//}
	//fmt.Println("result:", result)
	//result, err := rdb.SAdd(c, "why", "asd", "qwe").Result()
	//time2 := time.Minute * 10
	//result, err := rdb.Expire(c, "why1", time2).Result()
	//if err != nil {
	//	fmt.Println("rdb.HSet err:", err)
	//}
	//fmt.Println("result:", result)
	//userIDs := []int64{7093048250676019200, 12}
	//var userInfos []*User
	//a := Mdb.Model(&User{})
	//b := a.Where("id IN ?", userIDs)
	//err := b.Find(&userInfos).Error
	//if err != nil {
	//	fmt.Println("rdb.HSet err:", err)
	//}
	//fmt.Printf("%#v\n", userInfos[0])
	//user := &User{
	//	BaseModel: BaseModel{
	//		ID:        7092746779023638545,
	//		CreatedAt: time.Now(),
	//		UpdatedAt: time.Now(),
	//	},
	//	Username:        "why",
	//	Password:        "123456",
	//	Avatar:          "https://q1.qlogo.cn/g?b=qq&nk=1780006511&s=640",
	//	FollowCount:     10,
	//	FollowerCount:   20,
	//	BackgroundImage: "https://q1.qlogo.cn/g?b=qq&nk=1780006511&s=640",
	//	IsFollow:        false,
	//	Signature:       "没有了不好意思",
	//	TotalFavorited:  30,
	//	WorkCount:       40,
	//	FavoriteCount:   50,
	//}
	//// 存数据库
	////err := Mdb.Create(user).Error
	////if err != nil {
	////	fmt.Println("dao.Group.Mdb.Create(user) err:", err)
	////}
	//key := CreateUserKey(7092746779023638545)
	//userMap := CreateMapUserInfo(user)
	//err := rdb.HSet(c, key, userMap).Err()
	//if err != nil {
	//	fmt.Println("rdb.HSet err:", err)
	//}
	//result, err := rdb.HGetAll(c, key).Result()
	//if err != nil {
	//	fmt.Println("rdb.HGetAll err:", err)
	//}
	//fmt.Println("result:", result)
	//userInfo, err := CreateUserInfo(result)
	//if err != nil {
	//	fmt.Println("CreateUserInfo2 err:", err)
	//}
	//fmt.Println("userInfo", userInfo)
}

func CreateMapUserInfo2(userInfo *User) map[string]interface{} {
	userStr, _ := json.Marshal(userInfo)
	userMap := make(map[string]interface{})
	_ = json.Unmarshal(userStr, &userMap)
	fmt.Println("userMap:", userMap)
	return userMap
}

func CreateUserInfo2(userMap map[string]string) (*User, error) {
	userStr, _ := json.Marshal(userMap)
	userInfo := new(User)
	err := json.Unmarshal(userStr, userInfo)
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("userInfo:", userInfo)
	return userInfo, err
}

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.30.134:6379",
		Password: "", // 密码
		DB:       2,  // 数据库
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		// Redis连接失败，进行相应处理
		fmt.Println("redis初始化失败！！！！！")
		panic(err)
	}
	return rdb
}
func InitMysql() {
	//配置MySQL连接参数
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "123456", "192.168.30.134", "3309", "toktik-user")
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
