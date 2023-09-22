package ratelimit

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed lua/slide_window.lua
var luaSlideWindow string

type RedisSlideWindowLimiter struct {
	client   redis.Cmdable
	prefix   string
	interval int64
	maxCnt   int64
}

func (r *RedisSlideWindowLimiter) Limit(ctx context.Context) (bool, error) {
	return r.client.Eval(ctx, luaFixWindow, []string{r.prefix}, r.interval, r.maxCnt, time.Now().UnixMilli()).Bool()
}

func (r *RedisSlideWindowLimiter) Close(ctx context.Context) error {
	return nil
}

func NewRedisSlideWindowLimiter(client redis.Cmdable, interval time.Duration, maxCnt int64, opts ...RedisSlideWindowLimiterOption) *RedisSlideWindowLimiter {
	ret := &RedisSlideWindowLimiter{
		client:   client,
		prefix:   "my:rateLimit",
		interval: interval.Milliseconds(),
		maxCnt:   maxCnt,
	}

	for _, opt := range opts {
		opt(ret)
	}
	return ret

}

type RedisSlideWindowLimiterOption func(*RedisSlideWindowLimiter)

func RedisSlideWindowLimiterWithPrefix(prefix string) RedisSlideWindowLimiterOption {
	return func(limiter *RedisSlideWindowLimiter) {
		limiter.prefix = prefix
	}
}
