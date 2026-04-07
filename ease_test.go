package animate

import (
	"context"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEase(t *testing.T) {
	var delta float64 = 0.01
	if runtime.GOARCH == "wasm" {
		// wasm is.. slower and less predictable
		delta = 0.15
	}

	t.Run("ease over the given duration", func(t *testing.T) {
		const duration = 10 * time.Millisecond

		var logs []float64

		pacer := NewPacer(context.Background(), 2*time.Millisecond)
		animation := Ease{
			Pacer:    pacer,
			Duration: duration,
			Tick: func(progress float64) {
				logs = append(logs, progress)
			},
		}

		animation.Run(context.Background())

		for i, log := range logs {
			if i == len(logs)-1 {
				assert.Equal(t, 1.0, logs[len(logs)-1], "should end at exactly 1.0")
			}

			expected := float64(i+1) / float64(len(logs))
			assert.InDelta(t, expected, log, delta, "tick %d should be close to expected progress", i+1)
		}
	})

	t.Run("stop the animation when the context is cancelled", func(t *testing.T) {
		const duration = 10 * time.Millisecond

		var logs []float64

		ctx, cancel := context.WithCancel(context.Background())
		pacer := NewPacer(context.Background(), 2*time.Millisecond)
		animation := Ease{
			Pacer:    pacer,
			Duration: duration,
			Tick: func(progress float64) {
				logs = append(logs, progress)

				if progress >= 0.5 {
					cancel()
				}
			},
		}

		animation.Run(ctx)

		assert.NotEqual(t, 1.0, logs[len(logs)-1], "should not end at 1.0 due to cancellation")
	})

	t.Run("handle zero duration as an infinite animation", func(t *testing.T) {
		var logs []float64

		ctx, cancel := context.WithCancel(context.Background())
		pacer := NewPacer(context.Background(), 2*time.Millisecond)
		animation := Ease{
			Pacer:    pacer,
			Duration: 0,
			Tick: func(progress float64) {
				logs = append(logs, progress)

				if len(logs) >= 3 {
					cancel()
				}
			},
		}

		animation.Run(ctx)

		for _, log := range logs {
			assert.Equal(t, 0.0, log, "should always be 0 for zero duration")
		}
	})

	t.Run("handle negative duration as an infinite animation", func(t *testing.T) {
		var logs []float64

		ctx, cancel := context.WithCancel(context.Background())
		pacer := NewPacer(context.Background(), 2*time.Millisecond)
		animation := Ease{
			Pacer:    pacer,
			Duration: -1 * time.Second,
			Tick: func(progress float64) {
				logs = append(logs, progress)

				if len(logs) >= 3 {
					cancel()
				}
			},
		}

		animation.Run(ctx)

		for _, log := range logs {
			assert.Equal(t, 0.0, log, "should always be 0 for negative duration")
		}
	})

	t.Run("eases using the provided easing function", func(t *testing.T) {
		const duration = 10 * time.Millisecond

		var logs []float64

		pacer := NewPacer(context.Background(), 2*time.Millisecond)
		animation := Ease{
			Pacer:    pacer,
			Duration: duration,
			Easing:   EaseInSine,
			Tick: func(progress float64) {
				logs = append(logs, progress)
			},
		}

		animation.Run(context.Background())

		for i, log := range logs {
			if i == len(logs)-1 {
				assert.Equal(t, 1.0, logs[len(logs)-1], "should end at exactly 1.0")
			}

			expected := EaseInSine(float64(i+1) / float64(len(logs)))
			assert.InDelta(t, expected, log, delta, "tick %d should be close to expected eased progress", i+1)
		}
	})

	t.Run("does not panic when Tick is nil", func(t *testing.T) {
		pacer := NewPacer(context.Background(), time.Millisecond)
		animation := Ease{
			Pacer:    pacer,
			Duration: 10 * time.Millisecond,
			Tick:     nil,
		}

		assert.NotPanics(t, func() {
			animation.Run(context.Background())
		})
	})

	t.Run("can pause and resume the animation", func(t *testing.T) {
		const duration = 10 * time.Millisecond

		var wg sync.WaitGroup
		var mu sync.Mutex
		var logs []float64

		pacer := NewPacer(context.Background(), 2*time.Millisecond)
		animation := Ease{
			Pacer:    pacer,
			Duration: duration,
			Tick: func(progress float64) {
				mu.Lock()
				logs = append(logs, progress)
				mu.Unlock()
			},
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		wg.Go(func() { animation.Run(ctx) })

		time.Sleep(3 * time.Millisecond)
		animation.Pause()

		mu.Lock()
		assert.True(t, animation.IsPaused(), "should report as paused")
		assert.False(t, animation.IsRunning(), "should report as not running while paused")
		assert.Less(t, logs[len(logs)-1], 1.0, "should not have completed before pausing")
		mu.Unlock()

		time.Sleep(3 * time.Millisecond)
		animation.Resume()

		wg.Wait()

		mu.Lock()
		assert.False(t, animation.IsPaused(), "should report as not paused")
		assert.False(t, animation.IsRunning(), "should report as not running after completion")
		assert.Equal(t, 1.0, logs[len(logs)-1], "should end at exactly 1.0")
		mu.Unlock()
	})

	t.Run("can stop the animation", func(t *testing.T) {
		const duration = 10 * time.Millisecond

		var wg sync.WaitGroup
		var mu sync.Mutex
		var logs []float64

		pacer := NewPacer(context.Background(), 2*time.Millisecond)
		animation := Ease{
			Pacer:    pacer,
			Duration: duration,
			Tick: func(progress float64) {
				mu.Lock()
				logs = append(logs, progress)
				mu.Unlock()
			},
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		wg.Go(func() { animation.Run(ctx) })

		time.Sleep(3 * time.Millisecond)
		animation.Stop()

		mu.Lock()
		assert.False(t, animation.IsRunning(), "should report as not running after stopping")
		assert.False(t, animation.IsPaused(), "should report as not paused after stopping")
		mu.Unlock()

		wg.Wait()

		mu.Lock()
		assert.Less(t, logs[len(logs)-1], 1.0, "should not have completed before stopping")
		mu.Unlock()
	})
}
