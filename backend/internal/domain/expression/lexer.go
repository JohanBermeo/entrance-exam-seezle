package expression

import (
	"strconv"
	"unicode"
)

// Lexer converts an input string into a stream of tokens.
type Lexer struct {
	input    string
	position int
	start    int
	tokens   []Token
}

// NewLexer creates a new lexer for the given input.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

// Tokenize returns all tokens from the input.
func (l *Lexer) Tokenize() ([]Token, error) {
	for l.position < len(l.input) {
		l.start = l.position
		if err := l.scanToken(); err != nil {
			return nil, err
		}
	}
	l.tokens = append(l.tokens, Token{Type: TokenEOF, Start: l.position, End: l.position})
	return l.tokens, nil
}

func (l *Lexer) scanToken() error {
	if l.position >= len(l.input) {
		return nil
	}

	ch := rune(l.input[l.position])

	switch {
	case unicode.IsSpace(ch):
		l.position++
		return nil
	case unicode.IsDigit(ch) || ch == '.':
		return l.scanNumber()
	case unicode.IsLetter(ch) || ch == '_':
		return l.scanIdentifier()
	}

	switch ch {
	case '+':
		l.position++
		l.addToken(TokenPlus)
	case '-':
		l.position++
		l.addToken(TokenMinus)
	case '*':
		l.position++
		l.addToken(TokenMultiply)
	case '/':
		l.position++
		l.addToken(TokenDivide)
	case '^':
		l.position++
		l.addToken(TokenPower)
	case '(':
		l.position++
		l.addToken(TokenLParen)
	case ')':
		l.position++
		l.addToken(TokenRParen)
	case ',':
		l.position++
		l.addToken(TokenComma)
	default:
		return newLexerError(l.start, "unexpected character: "+string(ch))
	}
	return nil
}

func (l *Lexer) scanNumber() error {
	hasDot := false
	for l.position < len(l.input) {
		ch := rune(l.input[l.position])
		if unicode.IsDigit(ch) {
			l.position++
		} else if ch == '.' && !hasDot {
			hasDot = true
			l.position++
		} else {
			break
		}
	}

	value := l.input[l.start:l.position]
	if _, err := strconv.ParseFloat(value, 64); err != nil {
		return newLexerError(l.start, "invalid number: "+value)
	}

	l.addTokenWithValue(TokenNumber, value)
	return nil
}

func (l *Lexer) scanIdentifier() error {
	for l.position < len(l.input) {
		ch := rune(l.input[l.position])
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' {
			l.position++
		} else {
			break
		}
	}

	value := l.input[l.start:l.position]
	l.addTokenWithValue(TokenIdentifier, value)
	return nil
}

func (l *Lexer) addToken(t TokenType) {
	l.tokens = append(l.tokens, Token{Type: t, Start: l.start, End: l.position})
}

func (l *Lexer) addTokenWithValue(t TokenType, value string) {
	l.tokens = append(l.tokens, Token{Type: t, Value: value, Start: l.start, End: l.position})
}

func newLexerError(pos int, msg string) *ParseError {
	return &ParseError{Position: pos, Message: msg}
}
