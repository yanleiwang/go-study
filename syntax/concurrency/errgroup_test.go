package concurrency

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"sync/atomic"
	"testing"
)

/*
WaitGroup 和 errgroup.Group 是很相似的，可以
认为 errgroup.Group 是对 WaitGroup 的封装。
• 首先需要引入 golang.org/x/sync 依赖
• errgroup.Group 会帮我们保持进行中任务计数
• 任何一个任务返回 error，Wait 方法就会返回error
注意即使某个任务返回error，  wait 还是会等待所有任务完成，在返回
大多数情况下，随便选择哪个都可以，差异不大。

使用errgroup 实现 某个任务错误 即返回
			实现 超时控制：
https://keepmoving.ren/golang/sync-errgroup/

*/

func TestErrgroup(t *testing.T) {
	eg := errgroup.Group{}
	result := int64(0)
	for i := 0; i < 10; i++ {
		delta := i
		eg.Go(func() error {
			atomic.AddInt64(&result, int64(delta))
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}

	fmt.Println(result)

}
