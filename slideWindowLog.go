package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type SlideWindowLog struct {
	cancel       context.CancelFunc
	ctx          context.Context
	rds          *redis.Client
	limits       int
	timeInterval int
	bucketMap    map[string]struct{}
}

func SlideWindowLogInit(rds *redis.Client, limits int, timeInterval int) *SlideWindowLog {
	s := SlideWindowLog{
		rds:          rds,
		limits:       limits,
		timeInterval: timeInterval,
		bucketMap:    make(map[string]struct{}),
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())
	return &s
}

func (s *SlideWindowLog) RateLimit(limitKey string) bool {
	NowTime := time.Now()
	endTime := NowTime.Add(-time.Duration(s.timeInterval) * time.Millisecond).UnixNano()
	s.rds.ZRemRangeByScore(s.ctx, limitKey, "0", fmt.Sprintf("%v", endTime))

	amount, _ := s.rds.ZCard(s.ctx, limitKey).Result()

	if amount >= int64(s.limits) {
		return false
	}

	s.rds.ZAdd(s.ctx, limitKey, &redis.Z{Score: float64(time.Now().UnixNano()), Member: time.Now().UnixNano()})
	return true
}

func (s *SlideWindowLog) EndRateLimit() {
	s.rds.FlushAll(s.ctx)
}
