package animate

import (
	"context"
	"sync"
)

type Animation interface {
	Run(ctx context.Context)
}

// Run executes the given animation A and blocks until it is complete.
func Run(animations ...Animation) {
	RunContext(context.Background(), animations...)
}

// RunAsync executes the given animation A in a new goroutine without blocking.
func RunAsync(animations ...Animation) {
	RunContextAsync(context.Background(), animations...)
}

func RunContext(ctx context.Context, animations ...Animation) {
	var wg sync.WaitGroup

	for _, a := range animations {
		wg.Go(func() { a.Run(ctx) })
	}

	wg.Wait()
}

func RunContextAsync(ctx context.Context, animations ...Animation) {
	go RunContext(ctx, animations...)
}
