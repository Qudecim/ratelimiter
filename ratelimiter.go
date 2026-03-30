package reatelimiter

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
)

var RealisationType = int8(0) // 0 - semaphore, 1 - channel

type RateLimiter struct {
	limiter limiterInternal
}

func NewRateLimiter(limit int64) *RateLimiter {
	if RealisationType == int8(0) {
		return &RateLimiter{newSemaphoreLimiter(limit)}
	} else if RealisationType == int8(1) {
		return &RateLimiter{newChanLimiter(limit)}
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

type limiterInternal interface {
	Acquire(ctx context.Context) error
	TryAcquire() bool
	Release()
}

type semaphoreLimiter struct {
	sem *semaphore.Weighted
}

func newSemaphoreLimiter(limit int64) *semaphoreLimiter {
	return &semaphoreLimiter{semaphore.NewWeighted(limit)}
}

func (s *semaphoreLimiter) Acquire(ctx context.Context) error {
	if err := s.sem.Acquire(ctx, 1); err != nil {
		return fmt.Errorf("unable to acquire semaphore: %w", err)
	}
	return nil
}

func (s *semaphoreLimiter) TryAcquire() bool {
	return s.sem.TryAcquire(1)
}

func (s *semaphoreLimiter) Release() {
	s.sem.Release(1)
}

type chanLimiter struct {
	ch chan struct{}
}

func newChanLimiter(limit int64) *chanLimiter {
	return &chanLimiter{make(chan struct{}, limit)}
}

func (c *chanLimiter) Acquire(_ context.Context) error {
	c.ch <- struct{}{}
	return nil
}

func (c *chanLimiter) TryAcquire() bool {
	select {
	case <-c.ch:
		return true
	default:
		return false
	}
}

func (c *chanLimiter) Release() {
	<-c.ch
}
