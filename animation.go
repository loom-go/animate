package animate

import "sync/atomic"

type animation struct {
	gen    atomic.Uint64
	paused atomic.Bool
}

func (a *animation) IsIdle() bool {
	return a.gen.Load() == 0
}

func (a *animation) IsRunning() bool {
	return a.gen.Load() != 0 && !a.paused.Load()
}

func (a *animation) IsPaused() bool {
	return a.gen.Load() != 0 && a.paused.Load()
}

func (a *animation) Pause() {
	a.paused.Store(true)
}

func (a *animation) Resume() {
	a.paused.Store(false)
}

func (a *animation) Stop() {
	a.gen.Store(0)
	a.paused.Store(false)
}

func (a *animation) isCurrent(gen uint64) bool {
	return a.gen.Load() == gen
}

func (a *animation) start() uint64 {
	gen := a.gen.Add(1)
	a.paused.Store(false)

	return gen
}

func (a *animation) end(gen uint64) {
	if a.gen.CompareAndSwap(gen, 0) {
		a.paused.Store(false)
	}
}

func or[T comparable](v T, fallback T) T {
	var zero T
	if v == zero {
		return fallback
	}

	return v
}
