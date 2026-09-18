package calculation

import (
	"fmt"
)

// Graph represents a validated calculation DAG ready for execution.
type Graph struct {
	Operations []Operation
	Outputs    []string
	index      map[string]int
	adj        map[string][]string
	reverse    map[string][]string
}

// NewGraph builds and validates a calculation graph from operations and outputs.
func NewGraph(ops []Operation, outputs []string, knownOps map[string]int) (*Graph, error) {
	if len(ops) == 0 {
		return nil, NewDomainError(CodeInvalidInput, "", "at least one operation is required")
	}
	if len(outputs) == 0 {
		return nil, NewDomainError(CodeInvalidInput, "", "at least one output is required")
	}

	index := make(map[string]int, len(ops))
	for i, op := range ops {
		if _, exists := index[op.ID]; exists {
			return nil, NewDomainError(CodeInvalidInput, op.ID, "duplicate operation id")
		}
		if err := op.Validate(knownOps); err != nil {
			return nil, err
		}
		index[op.ID] = i
	}

	for _, outID := range outputs {
		if _, exists := index[outID]; !exists {
			return nil, NewDomainError(CodeUnknownReference, outID, "output references unknown operation")
		}
	}

	adj := make(map[string][]string, len(ops))
	reverse := make(map[string][]string, len(ops))
	for _, op := range ops {
		adj[op.ID] = nil
		reverse[op.ID] = nil
	}
	for _, op := range ops {
		for _, input := range op.Inputs {
			if input.IsRef() {
				refID := *input.Ref
				if _, exists := index[refID]; !exists {
					return nil, NewDomainError(CodeUnknownReference, op.ID, "references unknown operation "+refID)
				}
				if refID == op.ID {
					return nil, NewDomainError(CodeCycleDetected, op.ID, "operation cannot reference itself")
				}
				adj[refID] = append(adj[refID], op.ID)
				reverse[op.ID] = append(reverse[op.ID], refID)
			}
		}
	}

	if err := detectCycles(ops, adj); err != nil {
		return nil, err
	}

	return &Graph{
		Operations: ops,
		Outputs:    outputs,
		index:      index,
		adj:        adj,
		reverse:    reverse,
	}, nil
}

func detectCycles(ops []Operation, adj map[string][]string) error {
	color := make(map[string]int, len(ops))
	for _, op := range ops {
		color[op.ID] = 0
	}

	var dfs func(string) error
	dfs = func(id string) error {
		color[id] = 1
		for _, neighbor := range adj[id] {
			switch color[neighbor] {
			case 1:
				return NewDomainError(CodeCycleDetected, id, "cycle detected involving "+neighbor)
			case 0:
				if err := dfs(neighbor); err != nil {
					return err
				}
			}
		}
		color[id] = 2
		return nil
	}

	for _, op := range ops {
		if color[op.ID] == 0 {
			if err := dfs(op.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

// Dependencies returns the list of operation IDs that must complete before the given operation.
func (g *Graph) Dependencies(opID string) []string {
	return g.reverse[opID]
}

// Dependents returns the list of operation IDs that depend on the given operation.
func (g *Graph) Dependents(opID string) []string {
	return g.adj[opID]
}

// Operation returns the operation by ID, or nil if not found.
func (g *Graph) Operation(id string) *Operation {
	if idx, ok := g.index[id]; ok {
		return &g.Operations[idx]
	}
	return nil
}

// OperationIndex returns the index of the operation in the topological order.
func (g *Graph) OperationIndex(id string) (int, bool) {
	idx, ok := g.index[id]
	return idx, ok
}

// Len returns the number of operations in the graph.
func (g *Graph) Len() int {
	return len(g.Operations)
}

// ReadyOperations returns operations that have no unmet dependencies (for initial scheduling).
func (g *Graph) ReadyOperations(completed map[string]bool) []*Operation {
	var ready []*Operation
	for i := range g.Operations {
		op := &g.Operations[i]
		if completed[op.ID] {
			continue
		}
		allMet := true
		for _, input := range op.Inputs {
			if input.IsRef() {
				if !completed[*input.Ref] {
					allMet = false
					break
				}
			}
		}
		if allMet {
			ready = append(ready, op)
		}
	}
	return ready
}

// String returns a string representation of the graph for debugging.
func (g *Graph) String() string {
	var result string
	for _, op := range g.Operations {
		result += fmt.Sprintf("  %s: %s(", op.ID, op.Op)
		for i, input := range op.Inputs {
			if i > 0 {
				result += ", "
			}
			if input.IsLiteral() {
				result += fmt.Sprintf("%.2f", *input.Value)
			} else {
				result += "ref(" + *input.Ref + ")"
			}
		}
		result += ")\n"
	}
	result += "Outputs: "
	for i, out := range g.Outputs {
		if i > 0 {
			result += ", "
		}
		result += out
	}
	return result
}
