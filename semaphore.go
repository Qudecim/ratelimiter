package reatelimiter

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
)

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
