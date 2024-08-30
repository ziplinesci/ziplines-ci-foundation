package graceful

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

// DefaultTimeout is the default timeout for shutdown.
const DefaultTimeout = 30 * time.Second

var (
	shutdown     chan struct{}
	routinesDone chan struct{}
	observers    ObserverPool
	ErrTimeout   = errors.New("timeout waiting for shutdown")
)

func init() {
	initGrace()
}

// initGrace initializes the channels and the observer pool.
func initGrace() {
	shutdown = make(chan struct{})
	routinesDone = make(chan struct{})
	observers = NewObserverPool()
}

// kill (no param) default sends syscall.SIGTERM
// kill -2 is syscall.SIGINT
// kill -15 is syscall.SIGTERM.
var defaultOsSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}

// NewShutdownObserver adds a shutdown observer (goroutine) to the wait list.
// It returns a channel for listening to the shutdown signal and a close function
// to be called when the routine is done.
func NewShutdownObserver() (<-chan struct{}, func()) {
	closer := observers.Add()
	return shutdown, closer
}

// HandleSignals waits for the given signals. Once received, any of these signals
// will send a shutdown signal to goroutines listening on Shutdown.
// It waits for all goroutines to finish within the timeout duration before exiting.
// It should be called in the main goroutine to hold the process.
func HandleSignals(timeout time.Duration, signals ...os.Signal) error {
	return HandleSignalsWithContext(context.Background(), timeout, signals...)
}

// Shutdown sends a shutdown signal to goroutines listening on Shutdown.
// Goroutines suspended by calling HandleSignals or HandleSignalsWithContext will resume.
func Shutdown() error {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return fmt.Errorf("failed to find current process: %w", err)
	}
	if err = p.Signal(os.Interrupt); err != nil {
		return fmt.Errorf("failed to send interrupt signal: %w", err)
	}
	return nil
}

// HandleSignalsWithContext is the same as HandleSignals but with context support.
func HandleSignalsWithContext(ctx context.Context, timeout time.Duration, signals ...os.Signal) error {
	if len(signals) == 0 {
		signals = defaultOsSignals
	}
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	// Wait for interrupt signal to gracefully shutdown the server with
	event := observe(ctx, signals...)
	if event.Fired != nil {
		log.Info().Msgf("received %s(%#v)! shutting down", event.Fired.String(), event.Fired)
	}
	go triggerShutdown()
	log.Info().Msgf("waiting %d for services/ routines to finish\n", observers.Pending())
	select {
	case <-time.After(timeout):
		if observers.Pending() > 0 {
			log.Warn().Msgf("graceful shutdown: %d observers not closed: %w", observers.Pending(), ErrTimeout)
			return fmt.Errorf("graceful shutdown: %d observers not closed: %w", observers.Pending(), ErrTimeout)
		}
	case <-routinesDone:
	}
	log.Info().Msgf("all observers closed\n")
	if event.Fired == nil {
		return fmt.Errorf("graceful shutdown: %w", event.ContextErr)
	}
	return nil
}

func triggerShutdown() {
	close(shutdown)
	observers.Wait()
	close(routinesDone)
}

// observe waits for one of the provided signals or context cancellation.
func observe(ctx context.Context, signals ...os.Signal) *shutdownEvent {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, signals...)
	defer signal.Stop(sigCh)

	event := &shutdownEvent{}
	select {
	case fired := <-sigCh:
		event.Fired = fired
	case <-ctx.Done():
		event.ContextErr = ctx.Err()
	}
	return event
}

// shutdownEvent holds the details of the signal that triggered the shutdown or context error.
type shutdownEvent struct {
	Fired      os.Signal
	ContextErr error
}
