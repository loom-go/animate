package animate

import (
	"context"
	"math"
	"time"
)

type Spring struct {
	Tension  float64
	Friction float64
	Mass     float64

	From float64
	To   float64

	Precision float64

	Tick  func(value float64)
	Pacer *Pacer
}

func (s Spring) Run(ctx context.Context) {
	pacer := s.Pacer
	if pacer == nil {
		pacer = globalPacer
	}

	tension := s.Tension
	if tension == 0 {
		tension = 200
	}

	friction := s.Friction
	if friction == 0 {
		friction = 20
	}

	mass := s.Mass
	if mass == 0 {
		mass = 1
	}

	precision := s.Precision
	if precision == 0 {
		precision = 0.01
	}

	from := s.From
	to := s.To

	vel := 0.0
	prev := time.Now()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		settled := false
		pacer.Pace(func(now time.Time) {
			dt := now.Sub(prev).Seconds()
			prev = now

			force := -tension*(from-to) - friction*vel
			vel += (force / mass) * dt
			from += vel * dt

			s.Tick(from)

			if math.Abs(vel) < precision && math.Abs(from-to) < precision {
				from = to
				s.Tick(from)
				settled = true
			}
		})

		if settled {
			break
		}
	}
}
