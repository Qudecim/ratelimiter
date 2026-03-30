package reatelimiter

import "context"

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
