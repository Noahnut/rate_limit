package ratelimit

import (
	"github.com/go-redis/redis/v8"
)

type AlgorithmType int

type RateLimitAlgorithm interface {
	RateLimit(limitKey string) bool
	EndRateLimit()
}

type RateLimit struct {
	rds          *redis.Client
	limits       int
	timeInterval int
	algoritm     RateLimitAlgorithm
}

const TokenBucketType AlgorithmType = 0

// Can choose which rate limit algorithm want to use
func NewRateLimiter(rds *redis.Client) *RateLimit {
	r := RateLimit{
		rds: rds,
	}
	return &r
}

func (r *RateLimit) RateLimiterInit(limits int, timeInterval int, algoritmType AlgorithmType) {
	r.timeInterval, r.limits = timeInterval, limits

	switch algoritmType {
	case TokenBucketType:
		r.algoritm = TokenBucketInit(r.rds, r.limits, r.timeInterval)
	}

}

func (r *RateLimit) StopRateLimiter() {
	r.algoritm.EndRateLimit()
	r.rds.Close()
}

func (r *RateLimit) RateLimiterChecker(limitKey string) bool {
	return r.algoritm.RateLimit(limitKey)
}
