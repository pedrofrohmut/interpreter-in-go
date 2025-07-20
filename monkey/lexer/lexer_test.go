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
	if lx.getCh() != input[0] {
		t.Errorf("Expect lexer.ch to be %q, but got %q", input[0], lx.getCh())
	}
}

func TestNewLexerEmptyInput(t *testing.T) {
	input := ""
	lx := NewLexer(input)

	if lx.input != "" {
		t.Errorf("Expected lexer.input to be %q, but got %q", "", lx.input)
	}
	if lx.pos != 0 {
		t.Errorf("Expected lexer.position to be %q, but got %q", 0, lx.pos)
	}
}

func TestGetNextToken(t *testing.T) {
	input := "=+(){},;"
	lx := NewLexer(input)

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

	for i, tt := range tests {
		tok := lx.GetNextToken()

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

func TestGetNextToken2(t *testing.T) {
	input := `let five = 5;`
	lx := NewLexer(input)

	tests := []struct {
		expectedType 	token.TokenType
		expectedLiteral string
	} {
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
	}

	for i, tt := range tests {
		tok := lx.GetNextToken()

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

// TODO: TestGetNextToken3 with
// 	input := `
// 		let five = 5;
// 		let ten = 10;
//
// 		let add = fn(x, y) {
// 			x + y;
// 		};
//
// 		let result = add(five, ten);
// 	`
