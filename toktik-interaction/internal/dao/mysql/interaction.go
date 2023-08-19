package mysql

import (
	"context"
	"fmt"
	"github.com/Happy-Why/toktik-interaction/internal/model/auto"
	"gorm.io/gorm"
)

type InteractionDao struct {
	conn *GormConn
}

func NewInteractionDao() *InteractionDao {
	return &InteractionDao{
		conn: NewGormConn(),
	}
}

func (i *InteractionDao) IsRelationExist(c context.Context, myUserID, targetID int64) (uint, bool, error) {
	var count int64
	relation := new(auto.Relation)
	err := i.conn.Session(c).Model(&auto.Relation{}).
		Where("user_id = ? AND target_id = ?", myUserID, targetID).
		First(relation).Count(&count).Error
	return relation.ID, count > 0, err
}

func (i *InteractionDao) FollowUserAction(c context.Context, conn DbConn, relationInfo *auto.Relation) error {
	i.conn = conn.(*GormConn)
	return i.conn.Tx(c).Create(relationInfo).Error
}

func (i *InteractionDao) AddFollowCount(c context.Context, conn DbConn, userID uint) error {
	i.conn = conn.(*GormConn)
	return i.conn.Tx(c).Model(&auto.User{}).
		Where("id = ?", userID).
		Update("follow_count", gorm.Expr("follow_count + ?", 1)).Error
}

func (i *InteractionDao) AddFollowerCount(c context.Context, conn DbConn, userID uint) error {
	i.conn = conn.(*GormConn)
	return i.conn.Tx(c).Model(&auto.User{}).
		Where("id = ?", userID).
		Update("follower_count", gorm.Expr("follower_count + ?", 1)).Error
}

func (i *InteractionDao) CancelFollowUser(c context.Context, conn DbConn, relationInfo *auto.Relation) error {
	i.conn = conn.(*GormConn)
	return i.conn.Tx(c).Model(&auto.Relation{}).Unscoped().Delete(relationInfo).Error
}

func (i *InteractionDao) SubFollowCount(c context.Context, conn DbConn, userID uint) error {
	i.conn = conn.(*GormConn)
	return i.conn.Tx(c).Model(&auto.User{}).
		Where("id = ?", userID).
		Update("follow_count", gorm.Expr("follow_count - ?", 1)).Error
}

func (i *InteractionDao) SubFollowerCount(c context.Context, conn DbConn, userID uint) error {
	i.conn = conn.(*GormConn)
	return i.conn.Tx(c).Model(&auto.User{}).
		Where("id = ?", userID).
		Update("follower_count", gorm.Expr("follower_count - ?", 1)).Error
}

func (i *InteractionDao) GetFollowList(c context.Context, userID uint) ([]*auto.User, error) {
	followUsers := make([]*auto.User, 0)
	session := i.conn.Session(c)
	sql := fmt.Sprintf("SELECT u.* FROM relation r JOIN user u ON r.target_id = u.id WHERE r.user_id = ?;")
	raw := session.Raw(sql, userID)
	err := raw.Scan(&followUsers).Error
	return followUsers, err
}

func (i *InteractionDao) GetFansList(c context.Context, userID uint) ([]*auto.User, error) {
	fansUsers := make([]*auto.User, 0)
	session := i.conn.Session(c)
	sql := fmt.Sprintf("SELECT u.* FROM relation r JOIN user u ON r.user_id = u.id WHERE r.target_id = ?;")
	raw := session.Raw(sql, userID)
	err := raw.Scan(&fansUsers).Error
	return fansUsers, err
}

func (i *InteractionDao) GetFriendList(c context.Context, userID uint) ([]*auto.User, error) {
	friendUsers := make([]*auto.User, 0)
	session := i.conn.Session(c)
	sql := fmt.Sprintf(
		"SELECT u.* " +
			"FROM relation r2 " +
			"JOIN user u ON r2.user_id = u.id " +
			"WHERE r2.target_id = ? " +
			"AND r2.user_id IN " +
			"( SELECT r1.target_id FROM relation r1 WHERE user_id = ?);")
	raw := session.Raw(sql, userID, userID)
	err := raw.Scan(&friendUsers).Error
	return friendUsers, err
}
