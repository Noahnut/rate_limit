package ratelimit

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type FixWindowCounter struct {
	cancel       context.CancelFunc
	ctx          context.Context
	rds          *redis.Client
	windowTime   time.Time
	limits       int
	timeInterval int
}

func FixWindowCounterInit(rds *redis.Client, limits int, timeInterval int) *FixWindowCounter {
	f := FixWindowCounter{
		rds:          rds,
		limits:       limits,
		timeInterval: timeInterval,
		windowTime:   time.Now(),
	}

	f.ctx, f.cancel = context.WithCancel(context.Background())
	return &f
}

func (f *FixWindowCounter) RateLimit(limitKey string) bool {

	if _, exist := f.rds.Get(f.ctx, limitKey).Result(); exist == redis.Nil {
		f.rds.SetNX(f.ctx, limitKey, 1, time.Millisecond*time.Duration(f.timeInterval))
	} else {
		r, _ := f.rds.Incr(f.ctx, limitKey).Result()
		if r > int64(f.limits) {
			return false
		}
	}
	return true
}

func (f *FixWindowCounter) EndRateLimit() {
	f.rds.FlushAll(f.ctx)
}
