package animate

import (
	"context"
	"slices"
	"time"
)

// Timeline represents a sequence of keyframes to be executed at specific times.
type Timeline []Keyframe

type Keyframe struct {
	At time.Duration
	Do func(now, next time.Duration)
}

func (t Timeline) Run(ctx context.Context) {
	start := time.Now()

	ordered := slices.Clone(t)
	slices.SortFunc(ordered, func(a, b Keyframe) int {
		return int(a.At - b.At)
	})

	for i, keyframe := range ordered {
		select {
		case <-ctx.Done():
			return
		default:
		}

		var next time.Duration
		if i < len(t)-1 {
			next = t[i+1].At
		} else {
			next = keyframe.At
		}

		remaining := keyframe.At - time.Since(start)
		if remaining > 0 {
			time.Sleep(keyframe.At - time.Since(start))
		}

		keyframe.Do(time.Since(start), next)
	}
}
