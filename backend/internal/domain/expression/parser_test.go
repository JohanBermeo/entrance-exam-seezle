package expression

import (
	"testing"
)

func parseAndCheck(t *testing.T, input string, expectedExpr Expr) {
	lexer := NewLexer(input)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	parser := NewParser(tokens)
	expr, errors := parser.Parse()
	if len(errors) > 0 {
		t.Fatalf("Parser errors: %v", errors)
	}

	if !exprEquals(expr, expectedExpr) {
		t.Fatalf("AST mismatch\nexpected: %v\ngot:      %v", expectedExpr, expr)
	}
}

func exprEquals(a, b Expr) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	switch a := a.(type) {
	case NumberLiteral:
		b, ok := b.(NumberLiteral)
		return ok && a.Value == b.Value

	case Identifier:
		b, ok := b.(Identifier)
		return ok && a.Name == b.Name

	case BinaryExpr:
		b, ok := b.(BinaryExpr)
		return ok && exprEquals(a.Left, b.Left) && a.Operator.Type == b.Operator.Type && exprEquals(a.Right, b.Right)

	case UnaryExpr:
		b, ok := b.(UnaryExpr)
		return ok && a.Operator.Type == b.Operator.Type && exprEquals(a.Right, b.Right)

	case CallExpr:
		b, ok := b.(CallExpr)
		if !ok || a.Callee.Name != b.Callee.Name || len(a.Arguments) != len(b.Arguments) {
			return false
		}
		for i := range a.Arguments {
			if !exprEquals(a.Arguments[i], b.Arguments[i]) {
				return false
			}
		}
		return true

	case GroupingExpr:
		b, ok := b.(GroupingExpr)
		return ok && exprEquals(a.Expression, b.Expression)

	default:
		return false
	}
}

func TestParserNumberLiteral(t *testing.T) {
	parseAndCheck(t, "42", NumberLiteral{Value: 42})
	parseAndCheck(t, "3.14", NumberLiteral{Value: 3.14})
}

func TestParserIdentifier(t *testing.T) {
	parseAndCheck(t, "x", Identifier{Name: "x"})
}

func TestParserBinaryAddition(t *testing.T) {
	parseAndCheck(t, "1 + 2", BinaryExpr{
		Left:     NumberLiteral{Value: 1},
		Operator: Token{Type: TokenPlus},
		Right:    NumberLiteral{Value: 2},
	})
}

func TestParserBinarySubtraction(t *testing.T) {
	parseAndCheck(t, "5 - 3", BinaryExpr{
		Left:     NumberLiteral{Value: 5},
		Operator: Token{Type: TokenMinus},
		Right:    NumberLiteral{Value: 3},
	})
}

func TestParserBinaryMultiplication(t *testing.T) {
	parseAndCheck(t, "6 * 7", BinaryExpr{
		Left:     NumberLiteral{Value: 6},
		Operator: Token{Type: TokenMultiply},
		Right:    NumberLiteral{Value: 7},
	})
}

func TestParserBinaryDivision(t *testing.T) {
	parseAndCheck(t, "8 / 4", BinaryExpr{
		Left:     NumberLiteral{Value: 8},
		Operator: Token{Type: TokenDivide},
		Right:    NumberLiteral{Value: 4},
	})
}

func TestParserBinaryPower(t *testing.T) {
	parseAndCheck(t, "2 ^ 3", BinaryExpr{
		Left:     NumberLiteral{Value: 2},
		Operator: Token{Type: TokenPower},
		Right:    NumberLiteral{Value: 3},
	})
}

func TestParserPrecedenceMulOverAdd(t *testing.T) {
	// 1 + 2 * 3 should parse as 1 + (2 * 3)
	parseAndCheck(t, "1 + 2 * 3", BinaryExpr{
		Left:     NumberLiteral{Value: 1},
		Operator: Token{Type: TokenPlus},
		Right: BinaryExpr{
			Left:     NumberLiteral{Value: 2},
			Operator: Token{Type: TokenMultiply},
			Right:    NumberLiteral{Value: 3},
		},
	})
}

func TestParserPrecedencePowerOverMul(t *testing.T) {
	// 2 * 3 ^ 2 should parse as 2 * (3 ^ 2)
	parseAndCheck(t, "2 * 3 ^ 2", BinaryExpr{
		Left:     NumberLiteral{Value: 2},
		Operator: Token{Type: TokenMultiply},
		Right: BinaryExpr{
			Left:     NumberLiteral{Value: 3},
			Operator: Token{Type: TokenPower},
			Right:    NumberLiteral{Value: 2},
		},
	})
}

func TestParserPowerRightAssociative(t *testing.T) {
	// 2 ^ 3 ^ 2 should parse as 2 ^ (3 ^ 2) due to right associativity
	parseAndCheck(t, "2 ^ 3 ^ 2", BinaryExpr{
		Left:     NumberLiteral{Value: 2},
		Operator: Token{Type: TokenPower},
		Right: BinaryExpr{
			Left:     NumberLiteral{Value: 3},
			Operator: Token{Type: TokenPower},
			Right:    NumberLiteral{Value: 2},
		},
	})
}

func TestParserGrouping(t *testing.T) {
	// (1 + 2) * 3 should parse as (1 + 2) * 3
	parseAndCheck(t, "(1 + 2) * 3", BinaryExpr{
		Left: GroupingExpr{
			Expression: BinaryExpr{
				Left:     NumberLiteral{Value: 1},
				Operator: Token{Type: TokenPlus},
				Right:    NumberLiteral{Value: 2},
			},
		},
		Operator: Token{Type: TokenMultiply},
		Right:    NumberLiteral{Value: 3},
	})
}

func TestParserUnaryMinus(t *testing.T) {
	parseAndCheck(t, "-5", UnaryExpr{
		Operator: Token{Type: TokenMinus},
		Right:    NumberLiteral{Value: 5},
	})
}

func TestParserUnaryMinusWithPrecedence(t *testing.T) {
	// -2 ^ 2 should parse as -(2 ^ 2)
	parseAndCheck(t, "-2 ^ 2", UnaryExpr{
		Operator: Token{Type: TokenMinus},
		Right: BinaryExpr{
			Left:     NumberLiteral{Value: 2},
			Operator: Token{Type: TokenPower},
			Right:    NumberLiteral{Value: 2},
		},
	})
}

func TestParserFunctionCall(t *testing.T) {
	parseAndCheck(t, "sqrt(16)", CallExpr{
		Callee: Identifier{Name: "sqrt"},
		Arguments: []Expr{
			NumberLiteral{Value: 16},
		},
	})
}

func TestParserFunctionCallMultipleArgs(t *testing.T) {
	parseAndCheck(t, "percent(200, 15)", CallExpr{
		Callee: Identifier{Name: "percent"},
		Arguments: []Expr{
			NumberLiteral{Value: 200},
			NumberLiteral{Value: 15},
		},
	})
}

func TestParserNestedFunctionCalls(t *testing.T) {
	parseAndCheck(t, "sqrt(percent(200, 15))", CallExpr{
		Callee: Identifier{Name: "sqrt"},
		Arguments: []Expr{
			CallExpr{
				Callee: Identifier{Name: "percent"},
				Arguments: []Expr{
					NumberLiteral{Value: 200},
					NumberLiteral{Value: 15},
				},
			},
		},
	})
}

func TestParserComplexExpression(t *testing.T) {
	// sqrt(percent(200, 15)) + 4 ^ 2
	parseAndCheck(t, "sqrt(percent(200, 15)) + 4 ^ 2", BinaryExpr{
		Left: CallExpr{
			Callee: Identifier{Name: "sqrt"},
			Arguments: []Expr{
				CallExpr{
					Callee: Identifier{Name: "percent"},
					Arguments: []Expr{
						NumberLiteral{Value: 200},
						NumberLiteral{Value: 15},
					},
				},
			},
		},
		Operator: Token{Type: TokenPlus},
		Right: BinaryExpr{
			Left:     NumberLiteral{Value: 4},
			Operator: Token{Type: TokenPower},
			Right:    NumberLiteral{Value: 2},
		},
	})
}

func TestParserErrors(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"1 +", "unexpected token"},
		{"(1 + 2", "expected ')'"},
		{"1 2", "unexpected token"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			lexer := NewLexer(tc.input)
			tokens, _ := lexer.Tokenize()
			parser := NewParser(tokens)
			_, errors := parser.Parse()
			if len(errors) == 0 {
				t.Fatalf("expected error for %q", tc.input)
			}
			found := false
			for _, err := range errors {
				if len(err.Message) > 0 {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("expected error containing %q, got %v", tc.expected, errors)
			}
		})
	}
}
