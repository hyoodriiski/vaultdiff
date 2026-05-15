package diff

import (
	"errors"
	"testing"
	"time"
)

func TestWatch_DeliversResults(t *testing.T) {
	opts := WatchOptions{
		Enabled:  true,
		Interval: 1 * time.Millisecond,
		MaxRuns:  3,
	}

	call := 0
	fn := func() ([]Change, error) {
		call++
		return []Change{{Key: "k", Type: ChangeTypeModified}}, nil
	}

	stop := make(chan struct{})
	ch := Watch(opts, fn, stop)

	var results []WatchResult
	for r := range ch {
		results = append(results, r)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for i, r := range results {
		if r.Run != i+1 {
			t.Errorf("result %d: expected Run=%d, got %d", i, i+1, r.Run)
		}
		if len(r.Changes) != 1 {
			t.Errorf("result %d: expected 1 change, got %d", i, len(r.Changes))
		}
		if r.Err != nil {
			t.Errorf("result %d: unexpected error: %v", i, r.Err)
		}
	}
}

func TestWatch_PropagatesError(t *testing.T) {
	opts := WatchOptions{
		Interval: 1 * time.Millisecond,
		MaxRuns:  1,
	}
	expectedErr := errors.New("vault unavailable")
	fn := func() ([]Change, error) {
		return nil, expectedErr
	}

	stop := make(chan struct{})
	ch := Watch(opts, fn, stop)

	r := <-ch
	if r.Err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, r.Err)
	}
}

func TestWatch_StopsOnSignal(t *testing.T) {
	opts := WatchOptions{
		Interval: 10 * time.Millisecond,
		MaxRuns:  0, // unlimited
	}
	fn := func() ([]Change, error) {
		return nil, nil
	}

	stop := make(chan struct{})
	ch := Watch(opts, fn, stop)

	// receive one result then stop
	<-ch
	close(stop)

	// drain remaining
	count := 1
	for range ch {
		count++
	}

	if count > 3 {
		t.Errorf("expected watch to stop quickly, got %d results", count)
	}
}

func TestDefaultWatchOptions(t *testing.T) {
	opts := DefaultWatchOptions()
	if opts.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if opts.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", opts.Interval)
	}
	if opts.MaxRuns != 0 {
		t.Errorf("expected MaxRuns=0, got %d", opts.MaxRuns)
	}
}
