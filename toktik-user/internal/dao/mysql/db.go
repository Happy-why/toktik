package mysql

import (
	"context"
	"fmt"
	"github.com/Happy-Why/toktik-user/internal/dao"
	"github.com/Happy-Why/toktik-user/internal/global"
	"github.com/Happy-Why/toktik-user/internal/model/auto"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

func InitMysql() {
	//配置MySQL连接参数
	m := global.PvSettings.Mysql
	fmt.Println(m)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", m.Username, m.Password, m.Host, m.Port, m.DB)
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
	dao.Group.Mdb = DB
	_ = DB.AutoMigrate(&auto.Comment{}, &auto.Favorite{}, &auto.User{}, &auto.Relation{}, &auto.Video{})

}

func GetDB() *gorm.DB {
	return dao.Group.Mdb
}

type GormConn struct {
	db *gorm.DB
	tx *gorm.DB
}

func NewGormConn() *GormConn {
	return &GormConn{db: GetDB()}
}
func NewTran() *GormConn {
	return &GormConn{db: GetDB(), tx: GetDB()}
}
func (g *GormConn) Session(ctx context.Context) *gorm.DB {
	return g.db.Session(&gorm.Session{Context: ctx})
}

func (g *GormConn) Rollback() {
	g.tx.Rollback()
}
func (g *GormConn) Commit() {
	g.tx.Commit()
}

func (g *GormConn) Tx(ctx context.Context) *gorm.DB {
	return g.tx.WithContext(ctx)
}
