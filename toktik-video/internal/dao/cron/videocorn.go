package cron

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"toktik-common/timing_job"
	"toktik-video/internal/dao/mysql"
	"toktik-video/internal/dao/redis"
	"toktik-video/internal/model/auto"
)

func TimingJob() {
	go timing_job.StartMinuteJob(VideoInfoMoveToDB, 20, "videoInfo")
	go timing_job.StartMinuteJob(FavoriteRecordMoveToDB, 30, "favoriteRecord")
	select {}
}

func VideoInfoMoveToDB() {
	zap.L().Info("VideoInfos start move to DB !!!")
	ctx := context.Background()
	rClient := redis.GetRdbCache()
	dbConn := mysql.NewGormConn()
	// 从数据库取 video_ids
	videoIds := make([]int64, 0)
	err := dbConn.Session(ctx).Model(&auto.Video{}).Pluck("id", &videoIds).Error
	if err != nil {
		zap.L().Error("TimingJob VideoInfoMoveToDB GetVideoIds err:", zap.Error(err))
		return
	}
	// 去 redis中获取video_info
	for _, v := range videoIds {
		videoKey := auto.CreateVideoKey(uint(v))
		videoMap, err := rClient.HGetAll(ctx, videoKey)
		if err != nil {
			zap.L().Error("TimingJob VideoInfoMoveToDB rdb.HGetAll err:", zap.Error(err))
			return
		}
		videoInfo, err := auto.CreateVideoInfo(videoMap)
		if err != nil {
			zap.L().Error("TimingJob VideoInfoMoveToDB auto.CreateVideoInfo(videoMap) err:", zap.Error(err))
			return
		}
		// 更新数据库
		err = dbConn.Session(ctx).Model(&auto.Video{}).Where("id = ?", v).Updates(&videoInfo).Error
		if err != nil {
			zap.L().Error("TimingJob VideoInfoMoveToDB DB err:", zap.Error(err))
			return
		}
	}
}

func FavoriteRecordMoveToDB() {
	zap.L().Info("FavoriteRecords start move to DB !!!")
	ctx := context.Background()
	rclient := redis.GetRdbCache()
	dbConn := mysql.NewGormConn()
	// 遍历key = user_favorite::* ，获得 一堆 user_id,
	// 将 user_id 点赞的视频 video_id 拿到，判断 标志为 "1" or "2" ,1是点赞记录，2是取消点赞
	// 将点赞记录持久化到mysql中，然后删除缓存中 标志位 为2 的值
	// ① 获得 keys
	keys, err := getKeys(ctx, "user_favorite::*")
	if err != nil {
		zap.L().Error("FavoriteRecordMoveToDB getKeys err:", zap.Error(err))
		return
	}
	if len(keys) == 0 {
		return
	}
	for _, key := range keys {
		// ② 使用 key 去 redis 中查询 value
		videoIdsStr, err := rclient.SGetAll(ctx, key)
		if err != nil {
			zap.L().Error("FavoriteRecordMoveToDB rClient.SGetAll err:", zap.Error(err))
			return
		}
		// ③ 将 user_id 从 key 中 拆出来
		userIdStr := strings.SplitN(key, "::", 2)[1]
		userId, _ := strconv.ParseInt(userIdStr, 10, 64)

		for _, videoIdStr := range videoIdsStr {
			// ④ 将 value 拆解，分出 标志位 1，2 和 videoId
			value := strings.SplitN(videoIdStr, "+", 2)
			isFavorite := value[0]
			videoId, _ := strconv.ParseInt(value[1], 10, 64)
			// ⑤ 根据标志位，选择 删除点赞记录 还是 添加点赞记录 ，1为添加，2为删除
			if isFavorite == "1" {
				// 添加点赞记录
				session := dbConn.Session(ctx)
				sql := fmt.Sprintf("INSERT IGNORE INTO favorite (user_id, video_id) VALUES (?,?);")
				if err = session.Exec(sql, userId, videoId).Error; err != nil {
					zap.L().Error("FavoriteRecordMoveToDB createRecord err:", zap.Error(err))
					return
				}
			} else if isFavorite == "2" {
				// 删除点赞记录
				session := dbConn.Session(ctx)
				sql := fmt.Sprintf("DELETE FROM favorite WHERE user_id = ? AND video_id = ?;")
				if err = session.Exec(sql, userId, videoId).Error; err != nil {
					zap.L().Error("FavoriteRecordMoveToDB delDBRecord err:", zap.Error(err))
					return
				}
				// ⑥ 删除 缓存 中 标志 为 2 的点赞记录
				_, err = rclient.SDel(ctx, key, videoIdStr)
				if err != nil {
					zap.L().Error("FavoriteRecordMoveToDB delRedisRecord err:", zap.Error(err))
					return
				}
			}
		}
	}
}

func getKeys(ctx context.Context, keyPatten string) ([]string, error) {
	//keys, cursor, err := GetRedisHelper().Scan(ctx, 0, keyPatten, 10).Result()
	keys, err := redis.GetRdbCache().GetKeys(ctx, keyPatten)
	if err != nil {
		return nil, err
	}
	return keys, err
}
