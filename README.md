<h1 align="center"><code>animate</code></h1>

<p align="center">An embeddable UI animation engine in Go.</p>

```go
animate.Run(animate.Ease{
    Duration: time.Second * 2,
    Easing:   animate.EaseOutQuad,
    Tick: func(progress float64) {
        fmt.Printf("Progress: %v\n", progress)
    },
})
```

## Features

- Easing, spring and keyframe animations
