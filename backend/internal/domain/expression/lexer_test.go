package expression

import (
	"testing"
)

func TestLexerTokenizesNumbers(t *testing.T) {
	testCases := []struct {
		input    string
		expected []TokenType
	}{
		{"123", []TokenType{TokenNumber, TokenEOF}},
		{"42.5", []TokenType{TokenNumber, TokenEOF}},
		{"0.123", []TokenType{TokenNumber, TokenEOF}},
		{"1 2 3", []TokenType{TokenNumber, TokenNumber, TokenNumber, TokenEOF}},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			lexer := NewLexer(tc.input)
			tokens, err := lexer.Tokenize()
			if err != nil {
				t.Fatalf("Tokenize error: %v", err)
			}
			if len(tokens) != len(tc.expected) {
				t.Fatalf("expected %d tokens, got %d", len(tc.expected), len(tokens))
			}
			for i, expectedType := range tc.expected {
				if tokens[i].Type != expectedType {
					t.Fatalf("token %d: expected %v, got %v", i, expectedType, tokens[i].Type)
				}
			}
		})
	}
}

func TestLexerTokenizesIdentifiers(t *testing.T) {
	testCases := []struct {
		input    string
		expected []string
	}{
		{"sqrt", []string{"sqrt"}},
		{"percent", []string{"percent"}},
		{"add", []string{"add"}},
		{"my_var", []string{"my_var"}},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			lexer := NewLexer(tc.input)
			tokens, err := lexer.Tokenize()
			if err != nil {
				t.Fatalf("Tokenize error: %v", err)
			}
			if len(tokens) != len(tc.expected)+1 { // +1 for EOF
				t.Fatalf("expected %d tokens, got %d", len(tc.expected)+1, len(tokens))
			}
			for i, expected := range tc.expected {
				if tokens[i].Type != TokenIdentifier {
					t.Fatalf("token %d: expected identifier, got %v", i, tokens[i].Type)
				}
				if tokens[i].Value != expected {
					t.Fatalf("token %d: expected %q, got %q", i, expected, tokens[i].Value)
				}
			}
		})
	}
}

func TestLexerTokenizesOperators(t *testing.T) {
	testCases := []struct {
		input    string
		expected []TokenType
	}{
		{"+", []TokenType{TokenPlus, TokenEOF}},
		{"-", []TokenType{TokenMinus, TokenEOF}},
		{"*", []TokenType{TokenMultiply, TokenEOF}},
		{"/", []TokenType{TokenDivide, TokenEOF}},
		{"^", []TokenType{TokenPower, TokenEOF}},
		{"(", []TokenType{TokenLParen, TokenEOF}},
		{")", []TokenType{TokenRParen, TokenEOF}},
		{",", []TokenType{TokenComma, TokenEOF}},
		{"+-*/^(),", []TokenType{TokenPlus, TokenMinus, TokenMultiply, TokenDivide, TokenPower, TokenLParen, TokenRParen, TokenComma, TokenEOF}},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			lexer := NewLexer(tc.input)
			tokens, err := lexer.Tokenize()
			if err != nil {
				t.Fatalf("Tokenize error: %v", err)
			}
			if len(tokens) != len(tc.expected) {
				t.Fatalf("expected %d tokens, got %d", len(tc.expected), len(tokens))
			}
			for i, expectedType := range tc.expected {
				if tokens[i].Type != expectedType {
					t.Fatalf("token %d: expected %v, got %v", i, expectedType, tokens[i].Type)
				}
			}
		})
	}
}

func TestLexerRejectsInvalidCharacters(t *testing.T) {
	lexer := NewLexer("@")
	_, err := lexer.Tokenize()
	if err == nil {
		t.Fatal("expected error for invalid character")
	}
}

func TestLexerTokenizesComplexExpression(t *testing.T) {
	input := "sqrt(percent(200, 15)) + 4 ^ 2"
	lexer := NewLexer(input)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize error: %v", err)
	}

	expectedTypes := []TokenType{
		TokenIdentifier, // sqrt
		TokenLParen,
		TokenIdentifier, // percent
		TokenLParen,
		TokenNumber, // 200
		TokenComma,
		TokenNumber, // 15
		TokenRParen,
		TokenRParen,
		TokenPlus,
		TokenNumber, // 4
		TokenPower,
		TokenNumber, // 2
		TokenEOF,
	}

	if len(tokens) != len(expectedTypes) {
		t.Fatalf("expected %d tokens, got %d", len(expectedTypes), len(tokens))
	}
	for i, expectedType := range expectedTypes {
		if tokens[i].Type != expectedType {
			t.Fatalf("token %d: expected %v, got %v (value: %q)", i, expectedType, tokens[i].Type, tokens[i].Value)
		}
	}
}
