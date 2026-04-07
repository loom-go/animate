package animate

import (
	"context"
	"math"
	"time"
)

type Spring struct {
	animation

	Tension  float64
	Friction float64
	Mass     float64

	From float64
	To   float64

	Precision float64

	Tick  func(value float64)
	Pacer *Pacer
}

func (a *Spring) Run(ctx context.Context) {
	if !a.IsIdle() {
		return
	}

	gen := a.start()
	defer a.end(gen)

	vel := 0.0
	pos := a.From
	prev := time.Now()

	for {
		if ctx.Err() != nil || !a.isCurrent(gen) {
			break
		}

		to := a.To
		if to != pos {
			prev = time.Now()
		}

		pacer := or(a.Pacer, globalPacer)
		tension := or(a.Tension, 200)
		friction := or(a.Friction, 20)
		mass := or(a.Mass, 1)
		precision := or(a.Precision, 0.01)

		settled := false
		pacer.Pace(func(now time.Time) {
			if ctx.Err() != nil || !a.isCurrent(gen) || a.IsPaused() {
				return
			}

			dt := now.Sub(prev).Seconds()
			prev = now

			force := -tension*(pos-to) - friction*vel
			vel += (force / mass) * dt
			pos += vel * dt

			a.tick(pos)

			if math.Abs(vel) < precision && math.Abs(pos-to) < precision {
				a.tick(to)
				settled = true
			}
		})

		if settled {
			break
		}
	}
}

func (a *Spring) tick(value float64) {
	if a.Tick != nil {
		a.Tick(value)
	}
}
