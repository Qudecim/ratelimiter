package example

import (
	"context"
	"fmt"
	"sync"

	reatelimiter "github.com/qudecim/ratelimiter"
)

func main() {
	rl := reatelimiter.NewRateLimiter(5)
	ctx := context.Background()

	wg := sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			wg.Add(1)

			if err := rl.Acquire(ctx); err != nil {
				panic(err)
			}
			defer rl.Release()

			fmt.Println("acquire ok")
		}()
	}

	wg.Wait()
}
