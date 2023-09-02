package cron

import (
	"context"
	"go.uber.org/zap"
	"toktik-common/timing_job"
	"toktik-user/internal/cache"
	"toktik-user/internal/dao/mysql"
	"toktik-user/internal/model/auto"
)

func TimingJob() {
	go timing_job.StartMinuteJob(UserInfoMoveToDB, 30, "userInfo")
	select {}
}

func UserInfoMoveToDB() {
	ctx := context.Background()
	rdb := cache.GetRdbCache()
	dbConn := mysql.NewGormConn()
	// 从数据库取 user_ids
	userIds := make([]int64, 0)
	err := dbConn.Session(ctx).Model(&auto.User{}).Pluck("id", &userIds).Error
	if err != nil {
		zap.L().Error("TimingJob UserInfoMoveToDB GetUserIds err:", zap.Error(err))
		return
	}
	// 去 redis中获取user_count
	for _, v := range userIds {
		userCntKey := auto.CreateUserCountKey(uint(v))
		userCntMap, err := rdb.HGetAll(ctx, userCntKey)
		if err != nil {
			zap.L().Error("TimingJob UserInfoMoveToDB rdb.HGetAll err:", zap.Error(err))
			return
		}
		userCntInfo, err := auto.CreateUserCountInfo(userCntMap)
		if err != nil {
			zap.L().Error("TimingJob UserInfoMoveToDB auto.CreateUserCountInfo err:", zap.Error(err))
			return
		}
		// 更新数据库
		err = dbConn.Session(ctx).Model(&auto.UserCount{}).Where("user_id = ?", v).Updates(&userCntInfo).Error
		if err != nil {
			zap.L().Error("TimingJob UserInfoMoveToDB DB err:", zap.Error(err))
			return
		}
	}
}
