package limit

import (
	"context"
	"sort"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter 限流接口
type RateLimiter interface {
	Wait(ctx context.Context) error // 阻塞等待
	Limit() rate.Limit
}

type multiLimiter struct {
	limiters []RateLimiter
}

// Wait 阻塞等待直到获取令牌或者超时
func (m *multiLimiter) Wait(ctx context.Context) error {
	for _, l := range m.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Limit 返回当前限制速率
func (m *multiLimiter) Limit() rate.Limit {
	return m.limiters[0].Limit() // 直接返回限制最多的元素
}

// MultiLimiter 混合多个限流桶
func MultiLimiter(limiters ...RateLimiter) *multiLimiter {
	byLimit := func(i, j int) bool { return limiters[i].Limit() < limiters[j].Limit() } // 细粒度在前
	sort.Slice(limiters, byLimit)
	return &multiLimiter{limiters: limiters}
}

// Per 返回速率为 每duration,eventCount个请求
func Per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}
