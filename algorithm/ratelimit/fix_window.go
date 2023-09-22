package ratelimit

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// FixWindowLimiter 固定窗口
// 算法要点:
//
// + 一个时间段内(窗口), 只能执行n个请求
// + 窗口是固定的,  一个窗口结束了, 才能够启动新的窗口,  不会滑动
//
// 实际实现的时候,  没有必要开一个定时器去刷新窗口,  而是实现一种类似的效果:
//
// + 记录当前窗口开始的时间t1 和窗口大小 duration
// + 请求进来的时候,  假如当前时间t2  >  t1+ duration. 表明需要开始新的窗口
type FixWindowLimiter struct {
	timestamp int64 // 窗口起始时间
	duration  int64 // 窗口大小
	maxCnt    int64 // 在这个窗口内允许通过的最大请求数
	cnt       int64
	mutex     sync.Mutex
}

func NewFixWindowLimiter(duration time.Duration, maxCnt int64) *FixWindowLimiter {
	return &FixWindowLimiter{
		timestamp: time.Now().UnixMilli(),
		duration:  duration.Milliseconds(),
		maxCnt:    maxCnt,
		cnt:       0,
	}
}

// Limit 原子操作版
func (f *FixWindowLimiter) Limit(ctx context.Context) (bool, error) {
	now := time.Now().UnixMilli()
	timestamp := atomic.LoadInt64(&f.timestamp)
	// 需要开启新的窗口
	if timestamp+f.duration < now {
		if atomic.CompareAndSwapInt64(&f.timestamp, timestamp, now) {
			// 可以用store 而不必cas,  因为同时只会有一个go程进入
			atomic.StoreInt64(&f.cnt, 0)
			//atomic.CompareAndSwapInt64(&f.cnt, cnt, 0)
		}
	}
	cnt := atomic.AddInt64(&f.cnt, 1)
	if cnt >= f.maxCnt {
		// 触发瓶颈了， 需要限流
		return true, nil
	}
	return false, nil

}

// LimitMutex 加锁版
func (f *FixWindowLimiter) LimitMutex(ctx context.Context) (bool, error) {

	f.mutex.Lock()

	// 需要开启新的窗口
	now := time.Now().UnixNano()
	if f.timestamp+f.duration < now {
		f.cnt = 0
		f.timestamp = now
	}
	// 请求数满了
	if f.cnt >= f.maxCnt {
		f.mutex.Unlock()
		return true, nil
	}
	f.cnt++
	f.mutex.Unlock()
	return false, nil

}

func (f *FixWindowLimiter) Close(ctx context.Context) error {
	return nil
}
