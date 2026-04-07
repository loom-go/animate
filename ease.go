package animate

import (
	"context"
	"math"
	"time"
)

// Ease represents an animation thaht progresses from 0 to 1 over a duration of time.
type Ease struct {
	animation

	Duration time.Duration
	Easing   func(progress float64) float64
	Tick     func(progress float64)
	Pacer    *Pacer
}

func (a *Ease) Run(ctx context.Context) {
	if !a.IsIdle() {
		return
	}

	gen := a.start()
	defer a.end(gen)

	finite := a.Duration > 0
	start := time.Now()

	for {
		if ctx.Err() != nil || !a.isCurrent(gen) {
			break
		}

		pacer := or(a.Pacer, globalPacer)
		easing := a.Easing
		if easing == nil {
			easing = EaseLinear
		}

		settled := false
		pacer.Pace(func(now time.Time) {
			if ctx.Err() != nil || !a.isCurrent(gen) || a.IsPaused() {
				return
			}

			elapsed := max(0, now.Sub(start))

			if !finite {
				a.tick(0)
				return
			}

			elapsed = min(elapsed, a.Duration)
			progress := float64(elapsed) / float64(a.Duration)

			if finite && now.Sub(start) >= a.Duration {
				a.tick(1)
				settled = true
			} else {
				a.tick(easing(progress))
			}
		})

		if settled {
			break
		}
	}
}

func (a *Ease) tick(progress float64) {
	if a.Tick != nil {
		a.Tick(progress)
	}
}

func EaseLinear(progress float64) float64 {
	return progress
}

func EaseInSine(x float64) float64 {
	return 1 - math.Cos((x*math.Pi)/2)
}

func EaseOutSine(x float64) float64 {
	return math.Sin((x * math.Pi) / 2)
}

func EaseInOutSine(x float64) float64 {
	return -(math.Cos(math.Pi*x) - 1) / 2
}

func EaseInQuad(x float64) float64 {
	return x * x
}

func EaseOutQuad(x float64) float64 {
	return 1 - (1-x)*(1-x)
}

func EaseInOutQuad(x float64) float64 {
	if x < 0.5 {
		return 2 * x * x
	}
	return 1 - math.Pow(-2*x+2, 2)/2
}

func EaseInCubic(x float64) float64 {
	return x * x * x
}

func EaseOutCubic(x float64) float64 {
	return 1 - math.Pow(1-x, 3)
}

func EaseInOutCubic(x float64) float64 {
	if x < 0.5 {
		return 4 * x * x * x
	}
	return 1 - math.Pow(-2*x+2, 3)/2
}

func EaseInQuart(x float64) float64 {
	return x * x * x * x
}

func EaseOutQuart(x float64) float64 {
	return 1 - math.Pow(1-x, 4)
}

func EaseInOutQuart(x float64) float64 {
	if x < 0.5 {
		return 8 * x * x * x * x
	}
	return 1 - math.Pow(-2*x+2, 4)/2
}

func EaseInQuint(x float64) float64 {
	return x * x * x * x * x
}

func EaseOutQuint(x float64) float64 {
	return 1 - math.Pow(1-x, 5)
}

func EaseInOutQuint(x float64) float64 {
	if x < 0.5 {
		return 16 * x * x * x * x * x
	}
	return 1 - math.Pow(-2*x+2, 5)/2
}

func EaseInExpo(x float64) float64 {
	if x == 0 {
		return 0
	}
	return math.Pow(2, 10*x-10)
}

func EaseOutExpo(x float64) float64 {
	if x == 1 {
		return 1
	}
	return 1 - math.Pow(2, -10*x)
}

func EaseInOutExpo(x float64) float64 {
	switch {
	case x == 0:
		return 0
	case x == 1:
		return 1
	case x < 0.5:
		return math.Pow(2, 20*x-10) / 2
	default:
		return (2 - math.Pow(2, -20*x+10)) / 2
	}
}

func EaseInCirc(x float64) float64 {
	return 1 - math.Sqrt(1-math.Pow(x, 2))
}

func EaseOutCirc(x float64) float64 {
	return math.Sqrt(1 - math.Pow(x-1, 2))
}

func EaseInOutCirc(x float64) float64 {
	if x < 0.5 {
		return (1 - math.Sqrt(1-math.Pow(2*x, 2))) / 2
	}
	return (math.Sqrt(1-math.Pow(-2*x+2, 2)) + 1) / 2
}

func EaseInBack(x float64) float64 {
	const c1 = 1.70158
	return (c1+1)*x*x*x - c1*x*x
}

func EaseOutBack(x float64) float64 {
	const c1 = 1.70158
	return 1 + (c1+1)*math.Pow(x-1, 3) + c1*math.Pow(x-1, 2)
}

func EaseInOutBack(x float64) float64 {
	const c1 = 1.70158
	const c2 = c1 * 1.525
	if x < 0.5 {
		return (math.Pow(2*x, 2) * ((c2+1)*2*x - c2)) / 2
	}
	return (math.Pow(2*x-2, 2)*((c2+1)*(x*2-2)+c2) + 2) / 2
}

func EaseInElastic(x float64) float64 {
	switch {
	case x == 0:
		return 0
	case x == 1:
		return 1
	default:
		return -math.Pow(2, 10*x-10) * math.Sin((x*10-10.75)*((2*math.Pi)/3))
	}
}

func EaseOutElastic(x float64) float64 {
	switch {
	case x == 0:
		return 0
	case x == 1:
		return 1
	default:
		return math.Pow(2, -10*x)*math.Sin((x*10-0.75)*((2*math.Pi)/3)) + 1
	}
}

func EaseInOutElastic(x float64) float64 {
	const c5 = (2 * math.Pi) / 4.5
	switch {
	case x == 0:
		return 0
	case x == 1:
		return 1
	case x < 0.5:
		return -(math.Pow(2, 20*x-10) * math.Sin((20*x-11.125)*c5)) / 2
	default:
		return (math.Pow(2, -20*x+10)*math.Sin((20*x-11.125)*c5))/2 + 1
	}
}

func EaseOutBounce(x float64) float64 {
	const n1 = 7.5625
	const d1 = 2.75
	switch {
	case x < 1/d1:
		return n1 * x * x
	case x < 2/d1:
		x -= 1.5 / d1
		return n1*x*x + 0.75
	case x < 2.5/d1:
		x -= 2.25 / d1
		return n1*x*x + 0.9375
	default:
		x -= 2.625 / d1
		return n1*x*x + 0.984375
	}
}

func EaseInBounce(x float64) float64 {
	return 1 - EaseOutBounce(1-x)
}

func EaseInOutBounce(x float64) float64 {
	if x < 0.5 {
		return (1 - EaseOutBounce(1-2*x)) / 2
	}
	return (1 + EaseOutBounce(2*x-1)) / 2
}
