package animate

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeline(t *testing.T) {
	t.Run("runs each keyframe", func(t *testing.T) {
		var logs [][]time.Duration

		animation := Timeline{
			{At: 0, Do: func(now, next time.Duration) {
				logs = append(logs, []time.Duration{now, next})
			}},
			{At: 2 * time.Millisecond, Do: func(now, next time.Duration) {
				logs = append(logs, []time.Duration{now, next})
			}},
			{At: 4 * time.Millisecond, Do: func(now, next time.Duration) {
				logs = append(logs, []time.Duration{now, next})
			}},
		}

		animation.Run(context.Background())

		assert.Equal(t, 3, len(logs), "should have 3 keyframe executions")
		assert.InDelta(t, 0, logs[0][0].Seconds(), 0.01, "first keyframe should execute at approximately 0s")
		assert.InDelta(t, 2*time.Millisecond.Seconds(), logs[1][0].Seconds(), 0.01, "second keyframe should execute at approximately 2ms")
		assert.InDelta(t, 4*time.Millisecond.Seconds(), logs[2][0].Seconds(), 0.01, "third keyframe should execute at approximately 4ms")
	})

	t.Run("stops when the context is cancelled", func(t *testing.T) {
		var logs [][]time.Duration

		ctx, cancel := context.WithCancel(context.Background())
		animation := Timeline{
			{At: 0, Do: func(now, next time.Duration) {
				logs = append(logs, []time.Duration{now, next})
			}},
			{At: 2 * time.Millisecond, Do: func(now, next time.Duration) {
				logs = append(logs, []time.Duration{now, next})
				cancel()
			}},
			{At: 4 * time.Millisecond, Do: func(now, next time.Duration) {
				logs = append(logs, []time.Duration{now, next})
			}},
		}

		animation.Run(ctx)

		assert.Equal(t, 2, len(logs), "should have 2 keyframe executions before cancellation")
		assert.InDelta(t, 0, logs[0][0].Seconds(), 0.01, "first keyframe should execute at approximately 0s")
		assert.InDelta(t, 2*time.Millisecond.Seconds(), logs[1][0].Seconds(), 0.01, "second keyframe should execute at approximately 2ms")
	})

	t.Run("handles keyframes in any order", func(t *testing.T) {
		var logs [][]time.Duration

		animation := Timeline{
			{At: 4 * time.Millisecond, Do: func(now, next time.Duration) {
				logs = append(logs, []time.Duration{now, next})
			}},
			{At: 0, Do: func(now, next time.Duration) {
				logs = append(logs, []time.Duration{now, next})
			}},
			{At: 2 * time.Millisecond, Do: func(now, next time.Duration) {
				logs = append(logs, []time.Duration{now, next})
			}},
		}

		animation.Run(context.Background())

		assert.Equal(t, 3, len(logs), "should have 3 keyframe executions")
		assert.InDelta(t, 0, logs[0][0].Seconds(), 0.01, "first keyframe should execute at approximately 0s")
		assert.InDelta(t, 2*time.Millisecond.Seconds(), logs[1][0].Seconds(), 0.01, "second keyframe should execute at approximately 2ms")
		assert.InDelta(t, 4*time.Millisecond.Seconds(), logs[2][0].Seconds(), 0.01, "third keyframe should execute at approximately 4ms")
	})

	t.Run("handle keyframes with nil Do functions", func(t *testing.T) {
		var logs [][]time.Duration

		animation := Timeline{
			{At: 0, Do: func(now, next time.Duration) {
				logs = append(logs, []time.Duration{now, next})
			}},
			{At: 2 * time.Millisecond, Do: nil},
			{At: 4 * time.Millisecond, Do: func(now, next time.Duration) {
				logs = append(logs, []time.Duration{now, next})
			}},
		}

		animation.Run(context.Background())

		assert.Equal(t, 2, len(logs), "should have 2 keyframe executions")
		assert.InDelta(t, 0, logs[0][0].Seconds(), 0.01, "first keyframe should execute at approximately 0s")
		assert.InDelta(t, 4*time.Millisecond.Seconds(), logs[1][0].Seconds(), 0.01, "third keyframe should execute at approximately 4ms")
	})
}
