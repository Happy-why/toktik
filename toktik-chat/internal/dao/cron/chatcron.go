package cron

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
	"toktik-chat/internal/cache"
	"toktik-chat/internal/dao/mysql"
	"toktik-common/timing_job"
)

func TimingJob() {
	go timing_job.StartMinuteJob(ChatMessageMoveToDB, 60, "messageInfo")
	select {}
}

func ChatMessageMoveToDB() {
	zap.L().Info("ChatMessage start move to DB !!!")
	ctx := context.Background()
	rCache := cache.GetRdbCache()
	dbConn := mysql.NewGormConn()
	// ① 获得 keys chat::message::user_id+target_id
	keys, err := getKeys(ctx, "chat::message::*")
	if err != nil {
		zap.L().Error("ChatMessageMoveToDB getKeys err:", zap.Error(err))
		return
	}
	if len(keys) == 0 {
		return
	}
	for _, key := range keys {
		// ② 使用 key 去 redis 中查询 value
		messageList, err := rCache.ZGetRangeWithScores(ctx, key, 0, -1)
		if err != nil {
			zap.L().Error("ChatMessageMoveToDB rCache.ZGetRangeWithScores err:", zap.Error(err))
			return
		}
		for _, message := range messageList {
			// ③ 获取 created_at 并 将 user_id 和 target_id 从 value 中 拆出来
			createdTime := time.Unix(int64(message.Score), 0)
			str := strings.SplitN(message.Member.(string), "+", 3)
			userId, _ := strconv.ParseInt(str[0], 10, 64)
			targetId, _ := strconv.ParseInt(str[1], 10, 64)
			content := str[2]
			// ④ 构建 messageInfo 并 添加在数据库
			//messageInfo := &auto.Message{
			//	BaseModel: auto.BaseModel{CreatedAt: createdTime, UpdatedAt: createdTime},
			//	UserId:    uint(userId),
			//	ToUserId:  uint(targetId),
			//	Content:   content,
			//}
			session := dbConn.Session(ctx)
			sql := fmt.Sprintf("INSERT IGNORE INTO messages (created_at, updated_at, user_id, to_user_id, content) VALUES (?,?,?,?,?);")
			if err = session.Exec(sql, createdTime, createdTime, userId, targetId, content).Error; err != nil {
				zap.L().Error("FavoriteRecordMoveToDB createRecord err:", zap.Error(err))
				return
			}
			// ⑤ 将 redis 用于 持久化的 message 删掉
			_, err = rCache.ZDel(ctx, key, message.Member)
			if err != nil {
				zap.L().Error("ChatMessageMoveToDB DelMessageInfo err:", zap.Error(err))
				return
			}
		}

	}
}

func getKeys(ctx context.Context, keyPatten string) ([]string, error) {
	keys, err := cache.GetRdbCache().GetKeys(ctx, keyPatten)
	if err != nil {
		return nil, err
	}
	return keys, err
}
