// monkey/lexer/lexer_test.go

package lexer

import (
	"testing"
	"monkey/token"
)

func TestNewLexer(t *testing.T) {
	input := "bar"
	lx := NewLexer(input)

	if lx.input != input {
		t.Errorf("Expected lexer.input to be %q, but got %q", input, lx.input)
	}
	if lx.ch != input[0] {
		t.Errorf("Expect lexer.ch to be %q, but got %q", input[0], lx.ch)
	}
}

func TestNewLexerEmptyInput(t *testing.T) {
	input := ""
	lx := NewLexer(input)

	if lx.position != 0 {
		t.Errorf("Expected lexer.position to be %q, but got %q", 0, lx.position)
	}
	if lx.readPosition >= 1 {
		t.Errorf("Expected lexer.readPosition to be %q, but got %q", 0, lx.readPosition)
	}
	if lx.input != "" {
		t.Errorf("Expected lexer.input to be %q, but got %q", "", lx.input)
	}
	if lx.ch != 0 {
		t.Errorf("Expected lexer.ch to be %q, but got %q", 0, lx.ch)
	}
}

func TestNextToken(t *testing.T) {
	input := "=+(){},;"

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	} {
		{token.ASSIGN,    "="},
		{token.PLUS,      "+"},
		{token.LPAREN,    "("},
		{token.RPAREN,    ")"},
		{token.LBRACE,    "{"},
		{token.RBRACE,    "}"},
		{token.COMMA,     ","},
		{token.SEMICOLON, ";"},
		{token.EOF,       ""},
	}

	lx := NewLexer(input)

	for i, tt := range tests {
		tok := lx.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
