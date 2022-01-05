package ratelimit

import (
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestBuctetToken(t *testing.T) {
	rds := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,
	})

	defer rds.Close()
	rt := NewRateLimiter(rds)
	rt.RateLimiterInit(5, 100, TokenBucketType)
	defer rt.StopRateLimiter()

	for i := 0; i < 5; i++ {
		result := rt.RateLimiterChecker("5")
		if !result {
			t.Error("Result not expect")
		}
	}

	result := rt.RateLimiterChecker("5")

	if result {
		t.Error("Result not expect")
	}

	time.Sleep(100 * time.Millisecond)
	result = rt.RateLimiterChecker("5")

	if !result {
		t.Error("Result not expect")
	}
}

func TestBucketTokenMultiUser(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(3)

	rds := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,
	})

	rt := NewRateLimiter(rds)
	rt.RateLimiterInit(10, 100, TokenBucketType)
	defer rt.StopRateLimiter()

	go func() {

		for i := 0; i < 10; i++ {
			result := rt.RateLimiterChecker("127.0.0.1")
			if !result {
				t.Error("Result not expect")
			}
		}

		result := rt.RateLimiterChecker("127.0.0.1")

		if result {
			t.Error("Result not expect")
		}

		time.Sleep(100 * time.Millisecond)
		result = rt.RateLimiterChecker("127.0.0.1")

		if !result {
			t.Error("Result not expect")
		}

		wg.Done()
	}()

	go func() {

		for i := 0; i < 10; i++ {
			result := rt.RateLimiterChecker("127.0.0.2")
			if !result {
				t.Error("Result not expect")
			}
		}

		result := rt.RateLimiterChecker("127.0.0.2")

		if result {
			t.Error("Result not expect")
		}

		time.Sleep(100 * time.Millisecond)
		result = rt.RateLimiterChecker("127.0.0.2")

		if !result {
			t.Error("Result not expect")
		}

		wg.Done()
	}()

	go func() {

		for i := 0; i < 10; i++ {
			result := rt.RateLimiterChecker("127.0.0.3")
			if !result {
				t.Error("Result not expect")
			}
		}

		result := rt.RateLimiterChecker("127.0.0.3")

		if result {
			t.Error("Result not expect")
		}

		time.Sleep(100 * time.Millisecond)
		result = rt.RateLimiterChecker("127.0.0.3")

		if !result {
			t.Error("Result not expect")
		}

		wg.Done()
	}()

	wg.Wait()
}
