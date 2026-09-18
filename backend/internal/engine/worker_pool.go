package engine

import (
	"context"
)

// Pool bounds how many node executions run concurrently per request.
// Tasks are launched as goroutines that hold a semaphore slot while running.
type Pool struct {
	sem chan struct{}
}

// NewPool builds a pool with at most maxWorkers concurrent slots.
// Non-positive values collapse to a single slot (sequential execution).
func NewPool(maxWorkers int) *Pool {
	if maxWorkers < 1 {
		maxWorkers = 1
	}
	return &Pool{sem: make(chan struct{}, maxWorkers)}
}

// Capacity reports the maximum concurrent slots.
func (p *Pool) Capacity() int {
	return cap(p.sem)
}

// Run executes task holding a pool slot. It blocks until a slot is free or
// ctx is done; on cancellation it reports false and the task never runs.
func (p *Pool) Run(ctx context.Context, task func()) bool {
	select {
	case <-ctx.Done():
		return false
	case p.sem <- struct{}{}:
	}
	defer func() { <-p.sem }()
	task()
	return true
}
