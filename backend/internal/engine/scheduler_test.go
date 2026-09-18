package engine_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"back-calculator/internal/domain/calculation"
	"back-calculator/internal/domain/operators"
	"back-calculator/internal/engine"
)

func testRegistry() operators.Registry { return operators.NewRegistry() }

// sequentialReference evaluates the graph in declaration order. The test DAGs
// below only reference earlier nodes, so one pass suffices.
func sequentialReference(t *testing.T, g *calculation.Graph, registry operators.Registry) map[string]float64 {
	t.Helper()
	values := make(map[string]float64)
	for _, op := range g.Operations {
		var inputs []float64
		for _, in := range op.Inputs {
			if in.IsLiteral() {
				inputs = append(inputs, *in.Value)
				continue
			}
			v, ok := values[*in.Ref]
			if !ok {
				t.Fatalf("reference %s not computed yet", *in.Ref)
			}
			inputs = append(inputs, v)
		}
		v, err := registry.Evaluate(op.Op, inputs)
		if err != nil {
			t.Fatalf("Evaluate(%s) error = %v", op.Op, err)
		}
		values[op.ID] = v
	}
	return values
}

// randomAcyclicGraph builds a random DAG where node i only references nodes
// j < i, using total-safe operators (add/subtract/multiply over small ints).
func randomAcyclicGraph(rng *rand.Rand, n int) ([]calculation.Operation, []string) {
	ops := make([]calculation.Operation, 0, n)
	for i := 0; i < n; i++ {
		id := fmt.Sprintf("n%d", i)
		opNames := []string{"add", "subtract", "multiply"}
		op := opNames[rng.Intn(len(opNames))]
		inputs := make([]calculation.Input, 0, 2)
		for k := 0; k < 2; k++ {
			if i > 0 && rng.Intn(2) == 0 {
				inputs = append(inputs, calculation.Ref(fmt.Sprintf("n%d", rng.Intn(i))))
			} else {
				inputs = append(inputs, calculation.Literal(float64(rng.Intn(9)+1)))
			}
		}
		ops = append(ops, calculation.Operation{ID: id, Op: op, Inputs: inputs})
	}
	return ops, []string{fmt.Sprintf("n%d", n-1)}
}

func TestSchedulerMatchesSequentialOnRandomDAGs(t *testing.T) {
	registry := testRegistry()
	for seed := int64(0); seed < 25; seed++ {
		rng := rand.New(rand.NewSource(seed))
		n := 5 + rng.Intn(20)
		ops, outputs := randomAcyclicGraph(rng, n)
		g, err := calculation.NewGraph(ops, outputs, registry.ArityMap())
		if err != nil {
			t.Fatalf("seed %d: NewGraph error = %v", seed, err)
		}
		want := sequentialReference(t, g, registry)

		sched := engine.NewScheduler(registry, 4)
		got, err := sched.Execute(context.Background(), g)
		if err != nil {
			t.Fatalf("seed %d: Execute error = %v", seed, err)
		}
		for id, w := range want {
			res, ok := got.Get(id)
			if !ok {
				t.Fatalf("seed %d: missing result for %s", seed, id)
			}
			if res.Value != w {
				t.Fatalf("seed %d: result %s = %v, want %v", seed, id, res.Value, w)
			}
		}
	}
}

func TestSchedulerRespectsMaxWorkers(t *testing.T) {
	registry := testRegistry()
	const nodes = 16
	ops := make([]calculation.Operation, 0, nodes)
	for i := 0; i < nodes; i++ {
		ops = append(ops, calculation.Operation{
			ID:     fmt.Sprintf("n%d", i),
			Op:     "add",
			Inputs: []calculation.Input{calculation.Literal(1), calculation.Literal(2)},
		})
	}
	g, err := calculation.NewGraph(ops, []string{"n0"}, registry.ArityMap())
	if err != nil {
		t.Fatalf("NewGraph error = %v", err)
	}

	var current, maxSeen atomic.Int32
	exec := instrumentedExecutor{
		Registry: registry,
		before: func() {
			c := current.Add(1)
			for {
				m := maxSeen.Load()
				if c <= m || maxSeen.CompareAndSwap(m, c) {
					break
				}
			}
			time.Sleep(5 * time.Millisecond)
			current.Add(-1)
		},
	}

	sched := engine.Scheduler{Executor: exec, MaxWorkers: 3}
	results, err := sched.Execute(context.Background(), g)
	if err != nil {
		t.Fatalf("Execute error = %v", err)
	}
	if len(results) != nodes {
		t.Fatalf("got %d results, want %d", len(results), nodes)
	}
	if maxSeen.Load() > 3 {
		t.Fatalf("max concurrency = %d, want <= 3", maxSeen.Load())
	}
	if maxSeen.Load() < 2 {
		t.Fatalf("max concurrency = %d, want parallel execution (>1)", maxSeen.Load())
	}
}

func TestSchedulerFailsFastOnRequiredError(t *testing.T) {
	registry := testRegistry()
	ops := []calculation.Operation{
		{ID: "bad", Op: "divide", Inputs: []calculation.Input{calculation.Literal(1), calculation.Literal(0)}},
		{ID: "dep", Op: "add", Inputs: []calculation.Input{calculation.Ref("bad"), calculation.Literal(1)}},
		{ID: "ok", Op: "add", Inputs: []calculation.Input{calculation.Literal(2), calculation.Literal(3)}},
	}
	g, err := calculation.NewGraph(ops, []string{"dep"}, registry.ArityMap())
	if err != nil {
		t.Fatalf("NewGraph error = %v", err)
	}

	var mu sync.Mutex
	executed := make(map[string]bool)
	exec := instrumentedExecutor{
		Registry: registry,
		before:   nil,
		wrap: func(id string) {
			mu.Lock()
			executed[id] = true
			mu.Unlock()
		},
	}

	sched := engine.Scheduler{Executor: exec, MaxWorkers: 4}
	_, err = sched.Execute(context.Background(), g)
	if err == nil {
		t.Fatal("Execute error = nil, want division_by_zero")
	}
	var domainErr *calculation.DomainError
	if !errors.As(err, &domainErr) || domainErr.Code != calculation.CodeDivisionByZero {
		t.Fatalf("error = %v, want division_by_zero domain error", err)
	}
	mu.Lock()
	defer mu.Unlock()
	if executed["dep"] {
		t.Fatal("dependent of a failed node must not execute")
	}
}

func TestSchedulerCancellation(t *testing.T) {
	registry := testRegistry()
	ops := []calculation.Operation{
		{ID: "a", Op: "add", Inputs: []calculation.Input{calculation.Literal(1), calculation.Literal(2)}},
		{ID: "b", Op: "add", Inputs: []calculation.Input{calculation.Ref("a"), calculation.Literal(3)}},
	}
	g, err := calculation.NewGraph(ops, []string{"b"}, registry.ArityMap())
	if err != nil {
		t.Fatalf("NewGraph error = %v", err)
	}

	release := make(chan struct{})
	exec := instrumentedExecutor{
		Registry: registry,
		before: func() {
			<-release
		},
	}
	sched := engine.Scheduler{Executor: exec, MaxWorkers: 2}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := sched.Execute(ctx, g)
		done <- err
	}()
	time.Sleep(50 * time.Millisecond)
	cancel()
	close(release)

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("Execute error = nil, want cancellation error")
		}
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("error = %v, want context.Canceled", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Execute did not return after cancellation (deadlock?)")
	}
}

func TestSchedulerDiamondDeterminism(t *testing.T) {
	registry := testRegistry()
	ops := []calculation.Operation{
		{ID: "a", Op: "add", Inputs: []calculation.Input{calculation.Literal(10), calculation.Literal(5)}},
		{ID: "b", Op: "multiply", Inputs: []calculation.Input{calculation.Ref("a"), calculation.Literal(2)}},
		{ID: "c", Op: "subtract", Inputs: []calculation.Input{calculation.Ref("a"), calculation.Literal(3)}},
		{ID: "d", Op: "add", Inputs: []calculation.Input{calculation.Ref("b"), calculation.Ref("c")}},
	}
	g, err := calculation.NewGraph(ops, []string{"d", "b"}, registry.ArityMap())
	if err != nil {
		t.Fatalf("NewGraph error = %v", err)
	}
	sched := engine.NewScheduler(registry, 8)
	for i := 0; i < 20; i++ {
		results, err := sched.Execute(context.Background(), g)
		if err != nil {
			t.Fatalf("run %d: Execute error = %v", i, err)
		}
		if got := results.OutputValues([]string{"d", "b"}); len(got) != 2 || got[0] != 42 || got[1] != 30 {
			t.Fatalf("run %d: outputValues = %v, want [42 30]", i, got)
		}
	}
}

// instrumentedExecutor wraps RegistryExecutor with hooks for tests.
type instrumentedExecutor struct {
	Registry operators.Registry
	before   func()
	wrap     func(id string)
}

func (e instrumentedExecutor) Execute(op calculation.Operation, inputs []float64) (float64, error) {
	if e.wrap != nil {
		e.wrap(op.ID)
	}
	if e.before != nil {
		e.before()
	}
	return e.Registry.Evaluate(op.Op, inputs)
}
