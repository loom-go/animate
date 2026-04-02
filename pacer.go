package animate

import (
	"context"
	"sync"
	"time"
)

var globalPacer = NewPacer(context.Background(), time.Second/60)

// Pace paces a function at the global frame rate (60 FPS by default).
func Pace(tick func(time.Time)) {
	globalPacer.Pace(tick)
}

type frameRequest struct {
	tick func(now time.Time)
	done chan struct{}
}

type Pacer struct {
	ctx context.Context

	mu       sync.Mutex
	rate     time.Duration
	requests []*frameRequest
}

// NewPacer creates a new Pacer that paces frame updates at the given rate (frames per second).
// It can be given to animate.A to control the pacing of animations.
//
// By default, animations use a global FramePacer at 60 FPS.
func NewPacer(ctx context.Context, rate time.Duration) *Pacer {
	p := &Pacer{
		ctx:      ctx,
		rate:     rate,
		requests: make([]*frameRequest, 0),
	}

	go p.loop()
	return p
}

func (p *Pacer) loop() {
	defer p.drain()

	ticker := time.NewTicker(p.rate)
	defer ticker.Stop()

	for now := range ticker.C {
		select {
		case <-p.ctx.Done():
			return
		default:
		}

		p.mu.Lock()
		reqs := p.requests
		p.requests = nil
		p.mu.Unlock()

		if len(reqs) == 0 {
			continue
		}

		for _, req := range reqs {
			select {
			case <-p.ctx.Done():
			default:
				req.tick(now)
			}

			close(req.done)
		}
	}
}

func (p *Pacer) Pace(tick func(time.Time)) {
	req := &frameRequest{
		tick: tick,
		done: make(chan struct{}),
	}

	select {
	case <-p.ctx.Done():
		close(req.done)
	default:
		p.mu.Lock()
		p.requests = append(p.requests, req)
		p.mu.Unlock()
	}

	<-req.done
}

func (p *Pacer) drain() {
	p.mu.Lock()
	reqs := p.requests
	p.requests = nil
	p.mu.Unlock()

	for _, req := range reqs {
		close(req.done)
	}
}
