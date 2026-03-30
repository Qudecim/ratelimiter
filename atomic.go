package reatelimiter

import (
	"context"
	"sync"
	"sync/atomic"
)

type atomicLimiter struct {
	limit   int64
	counter int64
	cond    *sync.Cond
}

func newAtomicLimiter(limit int64) *atomicLimiter {
	return &atomicLimiter{limit, 0, sync.NewCond(&sync.Mutex{})}
}

func (a *atomicLimiter) Acquire(_ context.Context) error {
	if a.counter > a.limit {
		a.cond.Wait()
	}

	for !atomic.CompareAndSwapInt64(&a.counter, a.counter, a.counter+1) {
	}
	return nil
}

func (a *atomicLimiter) TryAcquire() bool {
	if a.counter > a.limit {
		return false
	}

	return atomic.CompareAndSwapInt64(&a.counter, a.counter, a.counter+1)
}

func (a *atomicLimiter) Release() {
	for !atomic.CompareAndSwapInt64(&a.counter, a.counter, a.counter-1) {
	}

	a.cond.Signal()
}
