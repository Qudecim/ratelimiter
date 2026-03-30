package reatelimiter

import (
	"context"
)

var RealisationType = int8(0) // 0 - semaphore, 1 - channel, 2 - atomic

type limiterInternal interface {
	Acquire(ctx context.Context) error
	TryAcquire() bool
	Release()
}

type RateLimiter struct {
	limiter limiterInternal
}

func NewRateLimiter(limit int64) *RateLimiter {
	if RealisationType == int8(0) {
		return &RateLimiter{newSemaphoreLimiter(limit)}
	} else if RealisationType == int8(1) {
		return &RateLimiter{newChanLimiter(limit)}
	} else if RealisationType == int8(2) {
		return &RateLimiter{newAtomicLimiter(limit)}
	}
	panic("realisationType must be 1 or 2")
}

func (r *RateLimiter) Acquire(ctx context.Context) error {
	return r.limiter.Acquire(ctx)
}

func (r *RateLimiter) TryAcquire() bool {
	return r.limiter.TryAcquire()
}

func (r *RateLimiter) Release() {
	r.limiter.Release()
}
