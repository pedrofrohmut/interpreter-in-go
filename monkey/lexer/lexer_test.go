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

type ExpectedToken struct {
    expectedType string
    expectedLiteral string
}

func checksForNextToken(lx *Lexer, t *testing.T, tests []ExpectedToken) {
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

func TestGetNextToken(t *testing.T) {
    input := "=+(){},;"
    lx := NewLexer(input)

    tests := []ExpectedToken {
        {token.Assign,    "="},
        {token.Plus,      "+"},
        {token.Lparen,    "("},
        {token.Rparen,    ")"},
        {token.Lbrace,    "{"},
        {token.Rbrace,    "}"},
        {token.Comma,     ","},
        {token.Semicolon, ";"},
        {token.Eof,       ""},
    }

    checksForNextToken(lx, t, tests)
}

func TestGetNextToken2(t *testing.T) {
    input := `let five = 5;`
    lx := NewLexer(input)

    tests := [] ExpectedToken {
        {token.Let, "let"},
        {token.Ident, "five"},
        {token.Assign, "="},
        {token.Int, "5"},
        {token.Semicolon, ";"},
    }

    checksForNextToken(lx, t, tests)
}

func TestGetNextToken3(t *testing.T) {
    input := `
        let five = 5;
        let ten = 10;

        let add = fn(x, y) {
            x + y;
        };

        let result = add(five, ten);
    `
    lx := NewLexer(input)

    tests := [] ExpectedToken {
        // Statement 1
        {token.Let, "let"},
        {token.Ident, "five"},
        {token.Assign, "="},
        {token.Int, "5"},
        {token.Semicolon, ";"},
        // Statement 2
        {token.Let, "let"},
        {token.Ident, "ten"},
        {token.Assign, "="},
        {token.Int, "10"},
        {token.Semicolon, ";"},
        // Statement 3
        {token.Let, "let"},
        {token.Ident, "add"},
        {token.Assign, "="},
        {token.Function, "fn"},
        {token.Lparen, "("},
        {token.Ident, "x"},
        {token.Comma, ","},
        {token.Ident, "y"},
        {token.Rparen, ")"},
        {token.Lbrace, "{"},
        {token.Ident, "x"},
        {token.Plus, "+"},
        {token.Ident, "y"},
        {token.Semicolon, ";"},
        {token.Rbrace, "}"},
        {token.Semicolon, ";"},
        // Statement 4
        {token.Let, "let"},
        {token.Ident, "result"},
        {token.Assign, "="},
        {token.Ident, "add"},
        {token.Lparen, "("},
        {token.Ident, "five"},
        {token.Comma, ","},
        {token.Ident, "ten"},
        {token.Rparen, ")"},
        {token.Semicolon, ";"},
    }

    checksForNextToken(lx, t, tests)
}

func TestGetNextToken4(t *testing.T) {
    input := `
        !-/*5;
        5 < 10 > 5;
    `
    lx := NewLexer(input)

    tests := [] ExpectedToken {
        // Statement 1
        {token.Bang, "!"},
        {token.Minus, "-"},
        {token.Slash, "/"},
        {token.Asterisk, "*"},
        {token.Int, "5"},
        {token.Semicolon, ";"},
        // Statement 2
        {token.Int, "5"},
        {token.Lt, "<"},
        {token.Int, "10"},
        {token.Gt, ">"},
        {token.Int, "5"},
        {token.Semicolon, ";"},
    }

    checksForNextToken(lx, t, tests)
}

func TestGetNextToken5(t *testing.T) {
    input := `
        if (5 < 10) {
            return true;
        } else {
            return false;
        }
    `
    lx := NewLexer(input)

    tests := []ExpectedToken {
        // Line 1
        {token.If, "if"},
        {token.Lparen, "("},
        {token.Int, "5"},
        {token.Lt, "<"},
        {token.Int, "10"},
        {token.Rparen, ")"},
        {token.Lbrace, "{"},
        // Line 2
        {token.Return, "return"},
        {token.True, "true"},
        {token.Semicolon, ";"},
        // Line 3
        {token.Rbrace, "}"},
        {token.Else, "else"},
        {token.Lbrace, "{"},
        // Line 4
        {token.Return, "return"},
        {token.False, "false"},
        {token.Semicolon, ";"},
        // Line 5
        {token.Rbrace, "}"},
    }

    checksForNextToken(lx, t, tests)
}

func TestNextToken6(t *testing.T) {
    input := `
        10 == 10;
        9 !=  10;
    `
    lx := NewLexer(input)

    tests := []ExpectedToken {
        // Statement 1
        {token.Int, "10"},
        {token.Eq, "=="},
        {token.Int, "10"},
        {token.Semicolon, ";"},
        // Statement 2
        {token.Int, "9"},
        {token.NotEq, "!="},
        {token.Int, "10"},
        {token.Semicolon, ";"},
    }

    checksForNextToken(lx, t, tests)
}

func TextNextTokenAllTokens(t *testing.T) {
    input := `foo 5 = + - ! * / < > == != , ; ( ) { } fn let true false if else return`
    lx := NewLexer(input)

    tests := []ExpectedToken {
        // Identifier + literals
        {token.Ident, "foo"},
        {token.Int, "5"},

        // Operators
        {token.Assign, "="},
        {token.Minus, "-"},
        {token.Bang, "!"},
        {token.Asterisk, "*"},
        {token.Slash, "/"},

        // Comparison
        {token.Lt, "<"},
        {token.Gt, ">"},
        {token.Eq, "=="},
        {token.NotEq, "!="},

        // Delimiters
        {token.Comma, ","},
        {token.Semicolon, ";"},

        // Grouping
        {token.Lparen, "("},
        {token.Rparen, ")"},
        {token.Lbrace, "{"},
        {token.Rbrace, "}"},

        // Keywords
        {token.Function, "fn"},
        {token.Let, "let"},
        {token.True, "true"},
        {token.False, "false"},
        {token.If, "if"},
        {token.Else, "else"},
        {token.Return, "return"},
    }

    checksForNextToken(lx, t, tests)
}
