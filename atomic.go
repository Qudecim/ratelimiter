package reatelimiter

import (
	"context"
	"sync"
	"sync/atomic"
)

type atomicLimiter struct {
	limit   int64
	counter atomic.Int64
	cond    *sync.Cond
}

func newAtomicLimiter(limit int64) *atomicLimiter {
	return &atomicLimiter{limit, atomic.Int64{}, sync.NewCond(&sync.Mutex{})}
}

func (a *atomicLimiter) Acquire(_ context.Context) error {
	n := a.counter.Load()

	if n > a.limit {
		a.cond.Wait()
	}

	for !atomic.CompareAndSwapInt64(&n, n, n+1) {
	}
	return nil
}

func (a *atomicLimiter) TryAcquire() bool {
	n := a.counter.Load()

	if n > a.limit {
		return false
	}

	return atomic.CompareAndSwapInt64(&n, n, n+1)
}

func (a *atomicLimiter) Release() {
	n := a.counter.Load()

	for !atomic.CompareAndSwapInt64(&n, n, n-1) {
	}

	a.cond.Signal()
}
