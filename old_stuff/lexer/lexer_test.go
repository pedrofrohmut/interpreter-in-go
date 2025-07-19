// monkey/lexer/lexer_test.go

package lexer

import (
	"testing"
	"monkey/token"
	"fmt"
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

func TestNextToken2(t *testing.T) {
	input := `let five = 5;`

	expected := []struct {
		ttype token.TokenType
		liter string
	} {
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
	}

	lx := NewLexer(input)

	for i, expec := range expected {
		nxt := lx.NextToken()

		fmt.Println(nxt.Type)
		fmt.Println(nxt.Literal)

		if expec.ttype != nxt.Type {
			t.Fatalf("tests[%d] - tokentype wrong. Expected = %q but got = %q", i, expec.ttype, nxt.Type)
		}
		if expec.liter != nxt.Literal {
			t.Fatalf("tests[%d] - literal wrong. Expected = %q but got = %q", i, expec.liter, nxt.Literal)
		}
	}
}

// func TextNextToken2(t *testing.T) {
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
//
// 	tests := []struct {
// 		expectedType token.TokenType
// 		expectedLiteral string
// 	} {
// 		// Statement 1
// 		{token.LET, "let"},
// 		{token.IDENT, "five"},
// 		{token.ASSIGN, "="},
// 		{token.INT, "5"},
// 		{token.SEMICOLON, ";"},
// 		// Statement 2
// 		{token.LET, "let"},
// 		{token.IDENT, "ten"},
// 		{token.ASSIGN, "="},
// 		{token.INT, "10"},
// 		{token.SEMICOLON, ";"},
// 		// Statement 3
// 		{token.LET, "let"},
// 		{token.IDENT, "add"},
// 		{token.ASSIGN, "="},
// 		{token.FUNCTION, "fn"},
// 		{token.LPAREN, "("},
// 		{token.IDENT, "x"},
// 		{token.COMMA, ","},
// 		{token.IDENT, "y"},
// 		{token.RPAREN, ")"},
// 		{token.LBRACE, "{"},
// 		{token.IDENT, "x"},
// 		{token.PLUS, "+"},
// 		{token.IDENT, "y"},
// 		{token.SEMICOLON, ";"},
// 		{token.RBRACE, "}"},
// 		{token.SEMICOLON, ";"},
// 		// Statement 4
// 		{token.LET, "let"},
// 		{token.IDENT, "result"},
// 		{token.ASSIGN, "="},
// 		{token.IDENT, "add"},
// 		{token.LPAREN, "("},
// 		{token.IDENT, "five"},
// 		{token.COMMA, ","},
// 		{token.IDENT, "ten"},
// 		{token.RPAREN, ")"},
// 		{token.SEMICOLON, ";"},
// 	}
// }
