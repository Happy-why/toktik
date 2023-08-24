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
	return i.conn.Tx(c).Model(&auto.Relation{}).
		Where("user_id = ? AND target_id = ?", relationInfo.UserId, relationInfo.TargetId).
		Unscoped().Delete(relationInfo).Error
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

func (i *InteractionDao) GetFollowIDs(c context.Context, userID uint) ([]int64, error) {
	var userIDs []int64
	session := i.conn.Session(c)
	sql := fmt.Sprintf("SELECT target_id FROM relation  WHERE user_id = ?;")
	raw := session.Raw(sql, userID)
	err := raw.Scan(&userIDs).Error
	return userIDs, err
}

func (i *InteractionDao) GetFansIDs(c context.Context, userID uint) ([]int64, error) {
	var userIDs []int64
	session := i.conn.Session(c)
	sql := fmt.Sprintf("SELECT user_id FROM relation  WHERE target_id = ?;")
	raw := session.Raw(sql, userID)
	err := raw.Scan(&userIDs).Error
	return userIDs, err
}

func (i *InteractionDao) GetFriendIDs(c context.Context, userID uint) ([]int64, error) {
	var userIDs []int64
	session := i.conn.Session(c)
	sql := fmt.Sprintf("SELECT user_id FROM relation  WHERE user_id IN  (SELECT target_id FROM relation  WHERE user_id = ?)AND target_id = ?;")
	raw := session.Raw(sql, userID, userID)
	err := raw.Scan(&userIDs).Error
	return userIDs, err
}

func (i *InteractionDao) IsFollowUser(c context.Context, myUserID, targetUserID int64) (bool, error) {
	var count int64
	err := i.conn.Session(c).Model(&auto.Relation{}).
		Where("user_id = ? AND target_id = ?", myUserID, targetUserID).Count(&count).Error
	return count > 0, err
}
