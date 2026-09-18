package main

import (
	"context"
	"sync"
)

// SendBulk fans a batch of notifications out across a bounded pool
// of goroutines. This is the piece that shows Go concurrency chops:
// a worker-pool pattern using a buffered channel as the job queue,
// so we never spawn an unbounded number of goroutines.

func (s *NotificationService) SendBulk(ctx context.Context, Notifications []*Notification, workers int) []error {
	jobs := make(chan *Notification)
	errs := make([]error, len(Notifications))

	var wg sync.WaitGroup
	// index protects mapping each error back to its notification
	// using a small wrapper struct sent through a results channel.
	type result struct {
		idx int
		err error
	}
	results := make(chan result, len(Notifications))
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := range jobs {
				// closures capture n correctly since n is the loop
				// variable of range over a channel (fresh each iter).
				idx := indexof(Notifications, n)
				err := s.Send(ctx, n)
				results <- result{idx: idx, err: err}
			}
		}()
	}

	go func() {
		for _, n := range Notifications {
			jobs <- n
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		errs[r.idx] = r.err
	}
	return errs
}

func indexof(list []*Notification, target *Notification) int {
	for i, n := range list {
		if n == target {
			return i
		}
	}
	return -1
}
