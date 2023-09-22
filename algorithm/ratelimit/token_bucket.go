package ratelimit

import (
	"context"
	"time"
)

var _ RateLimiter = (*TokenBucket)(nil)

// TokenBucket 令牌桶
// 算法要点:
// + 令牌每隔一段时间t生成,  放入容量为capacity的桶内
// + 每次请求从桶里拿走一个令牌
// + 拿到令牌的请求就会被处理
// + 没有拿到令牌的请求可以被:
//   - 直接被拒绝
//   - 阻塞直到拿到令牌或者超时
//   - 降级处理
type TokenBucket struct {
	tokens    chan struct{}
	closeChan chan struct{}
}

func NewTokenBucket(capacity int, duration time.Duration) *TokenBucket {

	tokens := make(chan struct{}, capacity)
	closeChan := make(chan struct{})
	go func() {
		ticker := time.NewTicker(duration)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				select {
				case tokens <- struct{}{}: // 可以放令牌
				default: // 桶满了， 无事发生
				}
			case <-closeChan:
				return

			}
		}
	}()
	return &TokenBucket{
		tokens:    tokens,
		closeChan: closeChan,
	}
}

func (b *TokenBucket) Limit(ctx context.Context) (bool, error) {
	select {
	case <-ctx.Done():
		// 超时
		return true, ctx.Err()
	case <-b.closeChan:
		// 令牌桶已经关闭
		return true, ErrLimiterClosed
	case <-b.tokens:
		// 拿到令牌了， 返回nil
		return false, nil
	}
}

func (b *TokenBucket) Close(ctx context.Context) error {
	close(b.closeChan)
	return nil
}
