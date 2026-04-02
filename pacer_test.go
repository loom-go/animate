package animate

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPacer(t *testing.T) {
	t.Run("runs at the correct pace", func(t *testing.T) {
		const frames = 5
		const duration = 20 * time.Millisecond

		pacer := NewPacer(context.Background(), duration)

		start := time.Now()
		for i := range frames {
			pacer.Pace(func(now time.Time) {
				elapsed := now.Sub(start)
				expected := duration * time.Duration(i+1)

				assert.InDelta(t, expected.Seconds(), elapsed.Seconds(), 0.01, "tick %d should be close to expected time", i+1)
			})
		}

		elapsed := time.Since(start)
		expected := duration * frames
		assert.InDelta(t, expected.Seconds(), elapsed.Seconds(), 0.01, "total elapsed time should be close to expected duration")
	})

	t.Run("handles context cancellation mid-pace", func(t *testing.T) {
		const frames = 2
		const duration = 10 * time.Millisecond

		var wg sync.WaitGroup
		var cancelled atomic.Bool

		ctx, cancel := context.WithCancel(context.Background())

		pacer := NewPacer(ctx, duration)

		for range frames {
			wg.Go(func() {
				pacer.Pace(func(now time.Time) {
					if cancelled.Load() {
						t.Errorf("tick should not run after cancellation")
					}

					cancel()
					cancelled.Store(true)
				})
			})
		}

		wg.Wait()
		cancel()
	})

	t.Run("handles pace after context cancellation", func(t *testing.T) {
		const duration = 10 * time.Millisecond

		ctx, cancel := context.WithCancel(context.Background())
		pacer := NewPacer(ctx, duration)

		cancel()

		pacer.Pace(func(now time.Time) {
			t.Error("Pace should not call the function after context cancellation")
		})
	})
}
