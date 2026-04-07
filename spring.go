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

func (s *Spring) Run(ctx context.Context) {
	vel := 0.0
	pos := s.From
	prev := time.Now()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		to := s.To
		if to != pos {
			prev = time.Now()
		}

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

		settled := false
		pacer.Pace(func(now time.Time) {
			dt := now.Sub(prev).Seconds()
			prev = now

			force := -tension*(pos-to) - friction*vel
			vel += (force / mass) * dt
			pos += vel * dt

			s.tick(pos)

			if math.Abs(vel) < precision && math.Abs(pos-to) < precision {
				s.tick(to)
				settled = true
			}
		})

		if settled {
			break
		}
	}
}

func (s *Spring) tick(value float64) {
	if s.Tick != nil {
		s.Tick(value)
	}
}
