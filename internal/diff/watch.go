package diff

import (
	"time"
)

// WatchOptions controls periodic re-diffing behavior.
type WatchOptions struct {
	Enabled  bool
	Interval time.Duration
	MaxRuns  int // 0 means unlimited
}

// DefaultWatchOptions returns watch options with sensible defaults.
func DefaultWatchOptions() WatchOptions {
	return WatchOptions{
		Enabled:  false,
		Interval: 30 * time.Second,
		MaxRuns:  0,
	}
}

// WatchResult holds the outcome of a single watch iteration.
type WatchResult struct {
	Run       int
	Timestamp time.Time
	Changes   []Change
	Err       error
}

// DiffFunc is a function that produces a slice of Changes on each call.
type DiffFunc func() ([]Change, error)

// Watch repeatedly calls fn at the configured interval, sending results to the
// returned channel. The channel is closed when MaxRuns is reached or the
// provided stop channel is closed.
func Watch(opts WatchOptions, fn DiffFunc, stop <-chan struct{}) <-chan WatchResult {
	out := make(chan WatchResult)

	go func() {
		defer close(out)
		run := 0

		for {
			changes, err := fn()
			run++
			result := WatchResult{
				Run:       run,
				Timestamp: time.Now(),
				Changes:   changes,
				Err:       err,
			}

			select {
			case out <- result:
			case <-stop:
				return
			}

			if opts.MaxRuns > 0 && run >= opts.MaxRuns {
				return
			}

			select {
			case <-time.After(opts.Interval):
			case <-stop:
				return
			}
		}
	}()

	return out
}
