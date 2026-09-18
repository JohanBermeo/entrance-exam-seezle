package engine

import (
	"back-calculator/internal/domain/calculation"
)

// dependencyPlan is the frozen scheduling state derived from a validated Graph.
// indegree counts unmet dependencies per node; dependents lists the nodes
// unblocked when a given node completes.
type dependencyPlan struct {
	indegree   map[string]int
	dependents map[string][]string
	total      int
}

// newDependencyPlan builds the scheduling state. The graph is already validated
// (unique IDs, existing refs, no cycles), so every ref resolves to a known node.
func newDependencyPlan(g *calculation.Graph) *dependencyPlan {
	plan := &dependencyPlan{
		indegree:   make(map[string]int, g.Len()),
		dependents: make(map[string][]string, g.Len()),
		total:      g.Len(),
	}
	for _, op := range g.Operations {
		plan.indegree[op.ID] = 0
		plan.dependents[op.ID] = nil
	}
	for _, op := range g.Operations {
		for _, input := range op.Inputs {
			if input.IsRef() {
				plan.indegree[op.ID]++
				plan.dependents[*input.Ref] = append(plan.dependents[*input.Ref], op.ID)
			}
		}
	}
	return plan
}

// roots returns the nodes with no dependencies, in declaration order.
func (p *dependencyPlan) roots(g *calculation.Graph) []string {
	var roots []string
	for _, op := range g.Operations {
		if p.indegree[op.ID] == 0 {
			roots = append(roots, op.ID)
		}
	}
	return roots
}
