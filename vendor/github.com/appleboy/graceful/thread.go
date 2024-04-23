package graceful

import "sync"

// routineGroup represents a group of goroutines.
type routineGroup struct {
	waitGroup sync.WaitGroup
}

// newRoutineGroup creates a new routineGroup.
func newRoutineGroup() *routineGroup {
	return new(routineGroup)
}

// Run runs a function in a new goroutine.
func (g *routineGroup) Run(fn func()) {
	g.waitGroup.Add(1)

	go func() {
		defer g.waitGroup.Done()
		fn()
	}()
}

// Wait waits for all goroutines to finish.
func (g *routineGroup) Wait() {
	g.waitGroup.Wait()
}
