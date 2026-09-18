package expression

// Node is the interface for all AST nodes.
type Node interface {
	node()
}

// Expr is the interface for expression nodes.
type Expr interface {
	Node
	expr()
}

// Stmt is the interface for statement nodes (not used in v1).
type Stmt interface {
	Node
	stmt()
}

// NumberLiteral represents a numeric literal.
type NumberLiteral struct {
	Value float64
	Token Token
}

func (NumberLiteral) node() {}
func (NumberLiteral) expr() {}

// Identifier represents a variable or function name.
type Identifier struct {
	Name  string
	Token Token
}

func (Identifier) node() {}
func (Identifier) expr() {}

// BinaryExpr represents a binary operation (e.g., a + b).
type BinaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (BinaryExpr) node() {}
func (BinaryExpr) expr() {}

// UnaryExpr represents a unary operation (e.g., -a).
type UnaryExpr struct {
	Operator Token
	Right    Expr
}

func (UnaryExpr) node() {}
func (UnaryExpr) expr() {}

// CallExpr represents a function call (e.g., sqrt(4)).
type CallExpr struct {
	Callee    Identifier
	Arguments []Expr
	Token     Token
}

func (CallExpr) node() {}
func (CallExpr) expr() {}

// GroupingExpr represents a parenthesized expression.
type GroupingExpr struct {
	Expression Expr
	Token      Token
}

func (GroupingExpr) node() {}
func (GroupingExpr) expr() {}
