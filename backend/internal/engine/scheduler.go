package engine

import (
	"context"
	"errors"
	"sync"

	"back-calculator/internal/domain/calculation"
	"back-calculator/internal/domain/operators"
)

// Scheduler executes a validated calculation DAG with bounded fan-out and a
// single fan-in collector. Independent nodes run in parallel; nodes with
// dependencies run once their references complete. On the first failure the
// scheduler cancels remaining work (fail-fast) and reports the root error,
// so the HTTP layer maps it to 422 (domain rule) or 408/499 (deadline/cancel).
type Scheduler struct {
	// Executor runs a single node. Use NewScheduler for the production
	// registry-backed executor, or set a custom one in tests.
	Executor NodeExecutor
	// MaxWorkers bounds concurrent node executions per request.
	// Non-positive values collapse to sequential execution.
	MaxWorkers int
}

// NewScheduler builds a Scheduler backed by the operator registry.
func NewScheduler(registry operators.Registry, maxWorkers int) Scheduler {
	return Scheduler{
		Executor:   RegistryExecutor{Registry: registry},
		MaxWorkers: maxWorkers,
	}
}

// nodeOutcome is the fan-in message a worker sends when it finishes a node.
type nodeOutcome struct {
	id    string
	value float64
	err   error
}

// Execute runs the graph to completion or to the first error.
// The returned results may be partial when err != nil.
func (s Scheduler) Execute(ctx context.Context, g *calculation.Graph) (calculation.Results, error) {
	exec := s.Executor
	if exec == nil {
		return nil, errors.New("engine: scheduler requires an Executor (use NewScheduler)")
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	plan := newDependencyPlan(g)
	pool := NewPool(s.MaxWorkers)
	done := make(chan nodeOutcome, plan.total)

	mu := &sync.RWMutex{}
	results := make(calculation.Results, plan.total)
	get := func(id string) (calculation.Result, bool) {
		mu.RLock()
		defer mu.RUnlock()
		res, ok := results[id]
		return res, ok
	}

	ready := plan.roots(g)
	pending := 0
	completed := 0
	failed := make(map[string]bool, plan.total)
	var firstErr error

	// release marks id settled and unblocks dependents, cascading failures
	// iteratively (no recursion: explicit DAG depth is unbounded in v1).
	release := func(stack []string) {
		for len(stack) > 0 {
			id := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, dep := range plan.dependents[id] {
				plan.indegree[dep]--
				if failed[id] {
					failed[dep] = true
				}
				if plan.indegree[dep] == 0 {
					if failed[dep] || firstErr != nil {
						failed[dep] = true
						completed++
						stack = append(stack, dep)
					} else {
						ready = append(ready, dep)
					}
				}
			}
		}
	}

	launch := func(id string) {
		pending++
		go func() {
			op := g.Operation(id)
			ran := pool.Run(ctx, func() {
				inputs, err := resolveInputs(op, get)
				var value float64
				if err == nil {
					value, err = exec.Execute(*op, inputs)
				}
				done <- nodeOutcome{id: id, value: value, err: err}
			})
			if !ran {
				done <- nodeOutcome{id: id, err: ctx.Err()}
			}
		}()
	}

	for _, id := range ready {
		launch(id)
	}
	ready = nil

	for completed < plan.total {
		if pending == 0 {
			return copyResults(mu, results), errors.New("engine: scheduler deadlock, no runnable nodes")
		}
		select {
		case <-ctx.Done():
			return copyResults(mu, results), ctx.Err()
		case out := <-done:
			pending--
			completed++
			if out.err != nil {
				if firstErr == nil {
					firstErr = out.err
					if errors.Is(out.err, context.Canceled) && ctx.Err() != nil {
						firstErr = ctx.Err()
					}
				}
				failed[out.id] = true
				cancel()
			} else {
				mu.Lock()
				results[out.id] = calculation.Result{ID: out.id, Value: out.value}
				mu.Unlock()
			}
			release([]string{out.id})
			for _, id := range ready {
				launch(id)
			}
			ready = nil
		}
	}

	return copyResults(mu, results), firstErr
}

func copyResults(mu *sync.RWMutex, results calculation.Results) calculation.Results {
	mu.RLock()
	defer mu.RUnlock()
	cp := make(calculation.Results, len(results))
	for id, res := range results {
		cp[id] = res
	}
	return cp
}
