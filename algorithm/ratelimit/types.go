package ratelimit

import (
	"context"
	"errors"
)

var ErrLimiterClosed = errors.New("限流器已关闭")

type RateLimiter interface {
	// Limit 返回值 bool值 表示是否需要限流
	Limit(ctx context.Context) (bool, error)
	Close(ctx context.Context) error
}
