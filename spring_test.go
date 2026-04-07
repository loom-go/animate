package animate

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSpring(t *testing.T) {
	t.Run("spring from x to y", func(t *testing.T) {
		var logs []float64

		pacer := NewPacer(context.Background(), time.Millisecond)
		animation := Spring{
			Pacer:   pacer,
			Tension: 200, Friction: 20,
			From: 0, To: 1, Tick: func(value float64) {
				logs = append(logs, value)
			},
		}

		animation.Run(context.Background())

		assert.GreaterOrEqual(t, len(logs), 2, "should have at least 2 ticks")
		assert.InDelta(t, animation.From, logs[0], 0.01, "should start at the 'From' value")
		assert.InDelta(t, animation.To, logs[len(logs)-1], 0.01, "should end at the 'To' value")
	})

	t.Run("stop the animation when the context is cancelled", func(t *testing.T) {
		var logs []float64

		ctx, cancel := context.WithCancel(context.Background())
		pacer := NewPacer(context.Background(), time.Millisecond)
		animation := Spring{
			Pacer:   pacer,
			Tension: 200, Friction: 20,
			From: 0, To: 1,
			Tick: func(value float64) {
				logs = append(logs, value)

				if value >= 0.5 {
					cancel()
				}
			},
		}

		animation.Run(ctx)

		assert.GreaterOrEqual(t, len(logs), 2, "should have at least 2 ticks before cancellation")
		assert.InDelta(t, animation.From, logs[0], 0.01, "should start at the 'From' value")
		assert.NotEqual(t, animation.To, logs[len(logs)-1], 0.01, "should not reach the 'To' value after cancellation")
	})

	t.Run("respects precision", func(t *testing.T) {
		var logs []float64

		pacer := NewPacer(context.Background(), time.Millisecond)
		animation := Spring{
			Pacer:   pacer,
			Tension: 200, Friction: 20,
			From: 0, To: 1,
			Precision: 0.001,
			Tick: func(value float64) {
				logs = append(logs, value)
			},
		}

		animation.Run(context.Background())

		assert.InDelta(t, 1.0, logs[len(logs)-1], 0.001)
	})

	t.Run("high friction does not overshoot", func(t *testing.T) {
		var logs []float64

		pacer := NewPacer(context.Background(), time.Millisecond)
		animation := Spring{
			Pacer:   pacer,
			Tension: 200, Friction: 80,
			From: 0, To: 1,
			Tick: func(v float64) { logs = append(logs, v) },
		}

		animation.Run(context.Background())

		for _, v := range logs {
			assert.LessOrEqual(t, v, 1.0+0.001, "should not overshoot the 'To' value with high friction")
		}
	})

	t.Run("low friction overshoots", func(t *testing.T) {
		var logs []float64

		pacer := NewPacer(context.Background(), time.Millisecond)
		animation := Spring{
			Pacer:   pacer,
			Tension: 200, Friction: 10,
			From: 0, To: 1,
			Tick: func(v float64) { logs = append(logs, v) },
		}

		animation.Run(context.Background())
	})

	t.Run("does not panic when Tick is nil", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		pacer := NewPacer(context.Background(), time.Millisecond)
		animation := Spring{
			Pacer:   pacer,
			Tension: 200, Friction: 20,
			From: 0, To: 1,
			Tick: nil,
		}

		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		assert.NotPanics(t, func() {
			animation.Run(ctx)
		})
	})

	t.Run("can pause and resume the animation", func(t *testing.T) {
		var wg sync.WaitGroup
		var mu sync.Mutex
		var logs []float64

		pacer := NewPacer(context.Background(), time.Millisecond)
		animation := Spring{
			Pacer:   pacer,
			Tension: 200, Friction: 20,
			From: 0, To: 1,
			Tick: func(value float64) {
				mu.Lock()
				logs = append(logs, value)
				mu.Unlock()
			},
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		wg.Go(func() { animation.Run(ctx) })

		time.Sleep(10 * time.Millisecond)
		animation.Pause()

		mu.Lock()
		assert.True(t, animation.IsPaused(), "should report as paused")
		assert.False(t, animation.IsRunning(), "should report as not running while paused")
		assert.Less(t, logs[len(logs)-1], 1.0, "should not have completed before pausing")
		mu.Unlock()

		time.Sleep(10 * time.Millisecond)
		animation.Resume()

		wg.Wait()

		mu.Lock()
		assert.False(t, animation.IsPaused(), "should report as not paused")
		assert.False(t, animation.IsRunning(), "should report as not running after completion")
		assert.InDelta(t, 1.0, logs[len(logs)-1], 0.01, "should end at approximately the 'To' value")
		mu.Unlock()
	})

	t.Run("can stop the animation", func(t *testing.T) {
		var wg sync.WaitGroup
		var mu sync.Mutex
		var logs []float64

		pacer := NewPacer(context.Background(), time.Millisecond)
		animation := Spring{
			Pacer:   pacer,
			Tension: 200, Friction: 20,
			From: 0, To: 1,
			Tick: func(value float64) {
				mu.Lock()
				logs = append(logs, value)
				mu.Unlock()
			},
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		wg.Go(func() { animation.Run(ctx) })

		time.Sleep(10 * time.Millisecond)
		animation.Stop()

		mu.Lock()
		assert.False(t, animation.IsRunning(), "should report as not running after stopping")
		assert.False(t, animation.IsPaused(), "should report as not paused after stopping")
		mu.Unlock()

		wg.Wait()

		mu.Lock()
		assert.False(t, animation.IsRunning(), "should report as not running after stopping")
		assert.False(t, animation.IsPaused(), "should report as not paused after stopping")
		assert.Less(t, logs[len(logs)-1], 1.0, "should not have completed after stopping")
		mu.Unlock()
	})
}
