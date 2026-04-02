package animate

type Animation interface {
	Run()
}

// Run executes the given animation A and blocks until it is complete.
func Run(a Animation) {
	a.Run()
}

// RunAsync executes the given animation A in a new goroutine without blocking.
func RunAsync(a Animation) {
	go a.Run()
}
