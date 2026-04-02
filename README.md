<h1 align="center"><code>animate</code></h1>

<p align="center">An embeddable UI animation engine in Go.</p>

```go
animate.Run(animate.Ease{
    Easing:   animate.EaseOutQuad,
    Tick: func(progress float64) {
        fmt.Printf("Progress: %v\n", progress)
    },
    Duration: time.Second * 2,
})
```

## Features

- Easings
- Springs
- Timelines and keyframes
- Precise frame pacing

## Usage

### Pacer

A pacer is responsible for running an animation smoothly at a given rate:

```go
// create a new 120FPS pacer that stops running when the context is cancelled
pacer := animate.NewPacer(context.Background(), time.Second/120)

// .Pace() schedules a function to be called in the next frame window
// and returns once the function has been called.
pacer.Pace(func(now time.Time) {
    // ...
})

for {
    pacer.Pace(func(now time.Time) {
        // runs on each frame (120FPS) until the loop is stopped
    })
}
```

`animate` keeps an internal global pacer running at 60FPS.
By default, if you don't specify a pacer to use for your animations, they will use the global pacer.

This global pacer can also be accesed with `animate.Pace()`:

```go
for {
    animate.Pace(func(now time.Time) {
        // runs at 60FPS
    })
}
```

### Easing

[Ease](https://easings.net/) a value from 0 to 1 over the given duration:

```go
// create a new Ease animation
a := animate.Ease{
    // Tick is called on each frame
    Tick: func(progress float64) {
        // `progress` goes from 0 to 1
        fmt.Printf("Progress: %v\n", progress)
    },

    // specify a duration
    Duration: time.Second * 2,

    // easing function to apply on `progress` (optional)
    Easing: animate.EaseOutQuad,

    // specify a custom pacer (optional)
    // else the global pacer is used (60FPS)
    Pacer: pacer,
}

// run the animation and stop if the context is cancelled
a.Run(context.Background())

// it can also be run with
animate.Run(a)
// or async
animate.RunAsync(a)
```

### Spring

Animate a value using a spring physics simulation:

```go
// create a new Spring animation
a := animate.Spring{
    // Tick is called on each frame
    Tick: func(value float64) {
        fmt.Printf("Value: %v\n", value)
    },

    // starting value
    From: 0,

    // target value
    To: 100,

    // spring physics parameters (optional)
    Tension:  200,
    Friction: 20,
    Mass:     1,

    // precision for stopping (optional)
    Precision: 0.01,

    // specify a custom pacer (optional)
    // else the global pacer is used (60FPS)
    Pacer: pacer,
}

// run the animation and stop if the context is cancelled
a.Run(context.Background())

// it can also be run with
animate.Run(a)
// or async
animate.RunAsync(a)
```

### Timeline

Run a series of keyframes (actions) at given times.
These keyframes can perform anything you want - from updating your UI, to running more animations (e.g. with `next - now` durations):

```go
animate.Run(animate.Timeline{
    {
        At: 0,
        Do: func(now, next time.Time) {
            // run an animation...
        },
    },
    {
        At: time.Second,
        Do: func(now, next time.Time) {
            // run an animation...
        },
    },
    {
        At: time.Second * 2,
        Do: func(now, next time.Time) {
            // run an animation...
        },
    },
})
```
