package ratelimit

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed lua/fix_window.lua
var luaFixWindow string

type RedisFixWindowLimiter struct {
	client   redis.Cmdable
	prefix   string
	interval int64
	maxCnt   int64
}

func (r *RedisFixWindowLimiter) Limit(ctx context.Context) (bool, error) {
	return r.client.Eval(ctx, luaFixWindow, []string{r.prefix}, r.interval, r.maxCnt).Bool()
}

func (r *RedisFixWindowLimiter) Close(ctx context.Context) error {
	return nil
}

func NewRediFixWindowLimiter(client redis.Cmdable, interval time.Duration, maxCnt int64, opts ...RedisFixWindowLimiterOption) *RedisFixWindowLimiter {
	ret := &RedisFixWindowLimiter{
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

type RedisFixWindowLimiterOption func(*RedisFixWindowLimiter)

func RedisFixWindowLimiterWithPrefix(prefix string) RedisFixWindowLimiterOption {
	return func(limiter *RedisFixWindowLimiter) {
		limiter.prefix = prefix
	}
}
