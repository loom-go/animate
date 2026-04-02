package animate

import "context"

type Animation interface {
	Run(ctx context.Context)
}

// Run executes the given animation A and blocks until it is complete.
func Run(a Animation) {
	RunCtx(context.Background(), a)
}

// RunAsync executes the given animation A in a new goroutine without blocking.
func RunAsync(a Animation) {
	RunCtxAsync(context.Background(), a)
}

func RunCtx(ctx context.Context, a Animation) {
	a.Run(ctx)
}

func RunCtxAsync(ctx context.Context, a Animation) {
	go RunCtx(ctx, a)
}
