package expression

import (
	"strconv"
)

// Precedence levels for the Pratt parser.
// Order from lowest to highest: AddSub < MulDiv < Unary < Power < Call < Primary
const (
	PrecedenceLowest int = iota
	PrecedenceAddSub
	PrecedenceMulDiv
	PrecedenceUnary
	PrecedencePower
	PrecedenceCall
	PrecedencePrimary
)

// Parser implements a Pratt parser for the expression language.
type Parser struct {
	tokens   []Token
	position int
	errors   []ParseError
}

// NewParser creates a new parser for the given tokens.
func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens}
}

// Parse parses the tokens and returns the AST root.
func (p *Parser) Parse() (Expr, []ParseError) {
	expr := p.parseExpression(PrecedenceLowest)
	if p.current().Type != TokenEOF {
		p.addError(p.current().Start, "unexpected token after expression: "+p.current().String())
	}
	return expr, p.errors
}

func (p *Parser) parseExpression(precedence int) Expr {
	left := p.parsePrimary()

	for precedence < p.getPrecedence(p.current()) {
		if p.current().Type == TokenEOF {
			break
		}

		switch p.current().Type {
		case TokenPlus, TokenMinus, TokenMultiply, TokenDivide, TokenPower:
			left = p.parseBinary(left)
		default:
			return left
		}
	}

	return left
}

func (p *Parser) parsePrimary() Expr {
	tok := p.current()

	switch tok.Type {
	case TokenNumber:
		p.advance()
		value, _ := strconv.ParseFloat(tok.Value, 64)
		return NumberLiteral{Value: value, Token: tok}

	case TokenIdentifier:
		p.advance()
		if p.current().Type == TokenLParen {
			return p.parseCall(Identifier{Name: tok.Value, Token: tok})
		}
		return Identifier{Name: tok.Value, Token: tok}

	case TokenMinus:
		p.advance()
		right := p.parseExpression(PrecedenceUnary)
		return UnaryExpr{Operator: tok, Right: right}

	case TokenPlus:
		p.advance()
		return p.parseExpression(PrecedenceUnary)

	case TokenLParen:
		p.advance()
		expr := p.parseExpression(PrecedenceLowest)
		if p.current().Type != TokenRParen {
			p.addError(p.current().Start, "expected ')'")
			return expr
		}
		parenTok := p.current()
		p.advance()
		return GroupingExpr{Expression: expr, Token: parenTok}

	default:
		p.addError(tok.Start, "unexpected token: "+tok.String())
		p.advance()
		return NumberLiteral{Value: 0, Token: tok}
	}
}

func (p *Parser) parseBinary(left Expr) Expr {
	op := p.current()
	p.advance()

	nextPrec := p.getPrecedence(op)
	// For right-associative operators (power), use lower precedence for right side
	if op.Type == TokenPower {
		nextPrec = PrecedencePower - 1
	}

	right := p.parseExpression(nextPrec)
	return BinaryExpr{Left: left, Operator: op, Right: right}
}

func (p *Parser) parseCall(name Identifier) Expr {
	tok := p.current()
	p.advance() // consume '('

	var args []Expr
	if p.current().Type != TokenRParen {
		for {
			args = append(args, p.parseExpression(PrecedenceLowest))
			if p.current().Type == TokenComma {
				p.advance()
			} else {
				break
			}
		}
	}

	if p.current().Type != TokenRParen {
		p.addError(p.current().Start, "expected ')' after function arguments")
	} else {
		p.advance()
	}

	return CallExpr{Callee: name, Arguments: args, Token: tok}
}

func (p *Parser) getPrecedence(tok Token) int {
	switch tok.Type {
	case TokenPlus, TokenMinus:
		return PrecedenceAddSub
	case TokenMultiply, TokenDivide:
		return PrecedenceMulDiv
	case TokenPower:
		return PrecedencePower
	default:
		return PrecedenceLowest
	}
}

func (p *Parser) current() Token {
	if p.position >= len(p.tokens) {
		return Token{Type: TokenEOF}
	}
	return p.tokens[p.position]
}

func (p *Parser) advance() {
	if p.position < len(p.tokens) {
		p.position++
	}
}

func (p *Parser) addError(pos int, msg string) {
	p.errors = append(p.errors, ParseError{Position: pos, Message: msg})
}

// ParseError represents a parsing error with position information.
type ParseError struct {
	Position int
	Message  string
}

func (e ParseError) Error() string {
	return "parse error at position " + strconv.Itoa(e.Position) + ": " + e.Message
}
