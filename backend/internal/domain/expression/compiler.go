package expression

import (
	"back-calculator/internal/domain/calculation"
	"back-calculator/internal/domain/operators"
	"fmt"
	"strconv"
)

// Compiler converts an AST into a calculation Graph.
type Compiler struct {
	registry    operators.Registry
	knownOps    map[string]int
	nodeCounter int
	operations  []calculation.Operation
	maxNodes    int
	maxDepth    int
}

// CompileResult holds either an operation ID or a literal value.
type CompileResult struct {
	OpID  string
	Value *float64
}

func (r CompileResult) IsLiteral() bool { return r.Value != nil }
func (r CompileResult) IsRef() bool     { return r.OpID != "" }

// NewCompiler creates a new compiler with the given operator registry.
func NewCompiler(registry operators.Registry) *Compiler {
	return &Compiler{
		registry: registry,
		knownOps: registry.ArityMap(),
		maxNodes: 100,
		maxDepth: 50,
	}
}

// SetLimits sets the maximum number of nodes and maximum AST depth.
func (c *Compiler) SetLimits(maxNodes, maxDepth int) {
	c.maxNodes = maxNodes
	c.maxDepth = maxDepth
}

// Compile compiles an expression AST into a calculation Graph.
// The final node is exposed under outputName (e.g. "result") so the
// HTTP layer can request it by the same name the client sent in outputs.
func (c *Compiler) Compile(expr Expr, outputName string) (*calculation.Graph, []ParseError) {
	c.nodeCounter = 0
	c.operations = nil
	if outputName == "" {
		outputName = "result"
	}

	var errors []ParseError
	result, err := c.compileExpr(expr, 0, &errors)
	if err != nil {
		return nil, errors
	}

	var outputID string
	if result.IsLiteral() {
		// Identity: literal + 0 preserves the value with a valid arity-2 add.
		outputID = outputName
		c.operations = append(c.operations, calculation.Operation{
			ID:     outputID,
			Op:     "add",
			Inputs: []calculation.Input{calculation.Literal(*result.Value), calculation.Literal(0)},
		})
	} else if result.OpID == outputName {
		outputID = result.OpID
	} else {
		// Identity alias: x + 0 preserves the value with a valid arity-2 add.
		outputID = outputName
		c.operations = append(c.operations, calculation.Operation{
			ID:     outputID,
			Op:     "add",
			Inputs: []calculation.Input{calculation.Ref(result.OpID), calculation.Literal(0)},
		})
	}

	graph, err := calculation.NewGraph(c.operations, []string{outputID}, c.knownOps)
	if err != nil {
		errors = append(errors, ParseError{Position: 0, Message: err.Error()})
		return nil, errors
	}

	return graph, errors
}

func (c *Compiler) compileExpr(expr Expr, depth int, errors *[]ParseError) (CompileResult, error) {
	if depth > c.maxDepth {
		*errors = append(*errors, ParseError{Position: 0, Message: "maximum AST depth exceeded"})
		return CompileResult{}, fmt.Errorf("max depth exceeded")
	}

	if len(c.operations) >= c.maxNodes {
		*errors = append(*errors, ParseError{Position: 0, Message: "maximum number of nodes exceeded"})
		return CompileResult{}, fmt.Errorf("max nodes exceeded")
	}

	switch e := expr.(type) {
	case NumberLiteral:
		return CompileResult{Value: &e.Value}, nil

	case Identifier:
		*errors = append(*errors, ParseError{Position: e.Token.Start, Message: "unbound identifier: " + e.Name})
		return CompileResult{}, fmt.Errorf("unbound identifier")

	case UnaryExpr:
		right, err := c.compileExpr(e.Right, depth+1, errors)
		if err != nil {
			return CompileResult{}, err
		}
		if e.Operator.Type == TokenMinus {
			if right.IsLiteral() {
				val := -*right.Value
				return CompileResult{Value: &val}, nil
			}
			id := c.newID("neg")
			c.operations = append(c.operations, calculation.Operation{
				ID:     id,
				Op:     "multiply",
				Inputs: []calculation.Input{calculation.Literal(-1), right.toInput()},
			})
			return CompileResult{OpID: id}, nil
		}
		return right, nil

	case BinaryExpr:
		left, err := c.compileExpr(e.Left, depth+1, errors)
		if err != nil {
			return CompileResult{}, err
		}
		right, err := c.compileExpr(e.Right, depth+1, errors)
		if err != nil {
			return CompileResult{}, err
		}

		opName := c.operatorName(e.Operator.Type)
		if opName == "" {
			*errors = append(*errors, ParseError{Position: e.Operator.Start, Message: "unknown operator: " + e.Operator.String()})
			return CompileResult{}, fmt.Errorf("unknown operator")
		}

		id := c.newID("bin")
		c.operations = append(c.operations, calculation.Operation{
			ID:     id,
			Op:     opName,
			Inputs: []calculation.Input{left.toInput(), right.toInput()},
		})
		return CompileResult{OpID: id}, nil

	case CallExpr:
		var argResults []CompileResult
		for _, arg := range e.Arguments {
			argResult, err := c.compileExpr(arg, depth+1, errors)
			if err != nil {
				return CompileResult{}, err
			}
			argResults = append(argResults, argResult)
		}

		opName := c.functionName(e.Callee.Name)
		if opName == "" {
			*errors = append(*errors, ParseError{Position: e.Token.Start, Message: "unknown function: " + e.Callee.Name})
			return CompileResult{}, fmt.Errorf("unknown function")
		}

		arity := c.knownOps[opName]
		if len(argResults) != arity {
			*errors = append(*errors, ParseError{Position: e.Token.Start, Message: fmt.Sprintf("function %s expects %d arguments, got %d", opName, arity, len(argResults))})
			return CompileResult{}, fmt.Errorf("wrong arity")
		}

		inputs := make([]calculation.Input, len(argResults))
		for i, arg := range argResults {
			inputs[i] = arg.toInput()
		}

		id := c.newID("call")
		c.operations = append(c.operations, calculation.Operation{
			ID:     id,
			Op:     opName,
			Inputs: inputs,
		})
		return CompileResult{OpID: id}, nil

	case GroupingExpr:
		return c.compileExpr(e.Expression, depth, errors)

	default:
		*errors = append(*errors, ParseError{Position: 0, Message: "unknown expression type"})
		return CompileResult{}, fmt.Errorf("unknown expression type")
	}
}

func (r CompileResult) toInput() calculation.Input {
	if r.IsLiteral() {
		return calculation.Literal(*r.Value)
	}
	return calculation.Ref(r.OpID)
}

func (c *Compiler) operatorName(tokType TokenType) string {
	switch tokType {
	case TokenPlus:
		return "add"
	case TokenMinus:
		return "subtract"
	case TokenMultiply:
		return "multiply"
	case TokenDivide:
		return "divide"
	case TokenPower:
		return "power"
	default:
		return ""
	}
}

func (c *Compiler) functionName(name string) string {
	switch name {
	case "sqrt":
		return "sqrt"
	case "percent":
		return "percent"
	default:
		return ""
	}
}

func (c *Compiler) newID(prefix string) string {
	c.nodeCounter++
	return prefix + strconv.Itoa(c.nodeCounter)
}
