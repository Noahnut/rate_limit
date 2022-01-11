package ratelimit

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type TokenBucket struct {
	cancel       context.CancelFunc
	ctx          context.Context
	rds          *redis.Client
	limits       int
	timeInterval int
	lock         sync.RWMutex
	bucketMap    map[string]byte
}

func TokenBucketInit(rds *redis.Client, limits int, timeInterval int) *TokenBucket {
	t := TokenBucket{
		rds:          rds,
		limits:       limits,
		timeInterval: timeInterval,
		bucketMap:    make(map[string]byte),
	}

	t.ctx, t.cancel = context.WithCancel(context.Background())
	go t.tokenAdder()
	return &t
}

func (t *TokenBucket) tokenAdder() {
	ticker := time.NewTicker(time.Duration(t.timeInterval) * time.Millisecond)
	defer ticker.Stop()
	select {
	case <-t.rds.Context().Done():
		return
	case <-ticker.C:
		for k := range t.bucketMap {
			_, err := t.rds.Get(t.ctx, k).Result()
			if err == redis.Nil {
				t.lock.Lock()
				delete(t.bucketMap, k)
				t.lock.Unlock()
			} else {
				t.rds.Incr(t.ctx, k)
			}
		}
	}
}

func (t *TokenBucket) RateLimit(limitKey string) bool {
	_, err := t.rds.Get(t.ctx, limitKey+"_count").Result()

	if err == redis.Nil {
		t.rds.Set(t.ctx, limitKey+"_count", t.limits, 10*time.Minute)
		t.lock.Lock()
		t.bucketMap[limitKey+"_count"] = 0
		t.lock.Unlock()
	} else {
		value, _ := t.rds.Get(t.ctx, limitKey+"_count").Result()
		requestLeft, _ := strconv.ParseInt(value, 10, 64)

		if requestLeft <= 0 {
			return false
		}
	}

	t.rds.Decr(t.ctx, limitKey+"_count")

	return true
}

func (t *TokenBucket) EndRateLimit() {

	for k := range t.bucketMap {
		t.rds.Del(t.ctx, k)
	}

	t.cancel()
}
