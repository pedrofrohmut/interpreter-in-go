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
    expectedType    token.TokenType
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

    checksForNextToken(lx, t, tests)
}

func TestGetNextToken2(t *testing.T) {
    input := `let five = 5;`
    lx := NewLexer(input)

    tests := [] ExpectedToken {
        {token.LET, "let"},
        {token.IDENT, "five"},
        {token.ASSIGN, "="},
        {token.INT, "5"},
        {token.SEMICOLON, ";"},
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
        {token.LET, "let"},
        {token.IDENT, "five"},
        {token.ASSIGN, "="},
        {token.INT, "5"},
        {token.SEMICOLON, ";"},
        // Statement 2
        {token.LET, "let"},
        {token.IDENT, "ten"},
        {token.ASSIGN, "="},
        {token.INT, "10"},
        {token.SEMICOLON, ";"},
        // Statement 3
        {token.LET, "let"},
        {token.IDENT, "add"},
        {token.ASSIGN, "="},
        {token.FUNCTION, "fn"},
        {token.LPAREN, "("},
        {token.IDENT, "x"},
        {token.COMMA, ","},
        {token.IDENT, "y"},
        {token.RPAREN, ")"},
        {token.LBRACE, "{"},
        {token.IDENT, "x"},
        {token.PLUS, "+"},
        {token.IDENT, "y"},
        {token.SEMICOLON, ";"},
        {token.RBRACE, "}"},
        {token.SEMICOLON, ";"},
        // Statement 4
        {token.LET, "let"},
        {token.IDENT, "result"},
        {token.ASSIGN, "="},
        {token.IDENT, "add"},
        {token.LPAREN, "("},
        {token.IDENT, "five"},
        {token.COMMA, ","},
        {token.IDENT, "ten"},
        {token.RPAREN, ")"},
        {token.SEMICOLON, ";"},
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
        {token.BANG, "!"},
        {token.MINUS, "-"},
        {token.SLASH, "/"},
        {token.ASTERISK, "*"},
        {token.INT, "5"},
        {token.SEMICOLON, ";"},
        // Statement 2
        {token.INT, "5"},
        {token.LT, "<"},
        {token.INT, "10"},
        {token.GT, ">"},
        {token.INT, "5"},
        {token.SEMICOLON, ";"},
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
        {token.IF, "if"},
        {token.LPAREN, "("},
        {token.INT, "5"},
        {token.LT, "<"},
        {token.INT, "10"},
        {token.RPAREN, ")"},
        {token.LBRACE, "{"},
        // Line 2
        {token.RETURN, "return"},
        {token.TRUE, "true"},
        {token.SEMICOLON, ";"},
        // Line 3
        {token.RBRACE, "}"},
        {token.ELSE, "else"},
        {token.LBRACE, "{"},
        // Line 4
        {token.RETURN, "return"},
        {token.FALSE, "false"},
        {token.SEMICOLON, ";"},
        // Line 5
        {token.RBRACE, "}"},
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
        {token.INT, "10"},
        {token.EQ, "=="},
        {token.INT, "10"},
        {token.SEMICOLON, ";"},
        // Statement 2
        {token.INT, "9"},
        {token.NOT_EQ, "!="},
        {token.INT, "10"},
        {token.SEMICOLON, ";"},
    }

    checksForNextToken(lx, t, tests)
}
