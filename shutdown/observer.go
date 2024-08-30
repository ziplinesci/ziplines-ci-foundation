package graceful

import (
	"sync"
	"sync/atomic"
)

// ObserverPool manages goroutines monitoring the graceful shutdown.
type observerPool struct {
	wg    *sync.WaitGroup
	count *atomic.Int64
	once  *sync.Once
}

// ObserverPool interface defines the public methods for managing observers.
type ObserverPool interface {
	Add() func()
	Pending() int
	Wait()
}

// NewObserverPool creates and returns a new ObserverPool.
func NewObserverPool() ObserverPool {
	return &observerPool{
		wg:    &sync.WaitGroup{},
		count: &atomic.Int64{},
		once:  &sync.Once{},
	}
}

func (o *observerPool) Add() func() {
	// Prevent adding new observers after the shutdown process has started.
	if o.count.Load() < 0 {
		return func() {} // return a no-op function
	}

	o.wg.Add(1)
	o.count.Add(1)
	return o.newCloser()
}

func (o *observerPool) newCloser() func() {
	var closed sync.Once

	return func() {
		closed.Do(func() {
			o.wg.Done()
			o.count.Add(-1)
		})
	}
}

func (o *observerPool) Pending() int {
	return int(o.count.Load())
}

func (o *observerPool) Wait() {
	o.wg.Wait()
}
