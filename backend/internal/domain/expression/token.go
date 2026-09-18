package expression

// TokenType represents the type of a token.
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenNumber
	TokenIdentifier
	TokenPlus
	TokenMinus
	TokenMultiply
	TokenDivide
	TokenPower
	TokenLParen
	TokenRParen
	TokenComma
)

// Token represents a single token from the lexer.
type Token struct {
	Type  TokenType
	Value string
	Start int
	End   int
}

// String returns a string representation of the token.
func (t Token) String() string {
	if t.Value != "" {
		return t.Value
	}
	switch t.Type {
	case TokenEOF:
		return "EOF"
	case TokenNumber:
		return "NUMBER"
	case TokenIdentifier:
		return "IDENTIFIER"
	case TokenPlus:
		return "+"
	case TokenMinus:
		return "-"
	case TokenMultiply:
		return "*"
	case TokenDivide:
		return "/"
	case TokenPower:
		return "^"
	case TokenLParen:
		return "("
	case TokenRParen:
		return ")"
	case TokenComma:
		return ","
	default:
		return "UNKNOWN"
	}
}

// IsOperator returns true if the token is a binary operator.
func (t Token) IsOperator() bool {
	return t.Type >= TokenPlus && t.Type <= TokenPower
}
