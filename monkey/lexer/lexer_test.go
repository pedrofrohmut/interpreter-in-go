// monkey/lexer/lexer_test.go

package lexer

import (
    "monkey/token"
    "testing"
)

func TestNewLexer(t *testing.T) {
    var input = "bar"
    var lexer = NewLexer(input)

    if lexer.input != input {
        t.Errorf("Expected lexer.input to be %q, but got %q", input, lexer.input)
    }
    if lexer.getCh() != input[0] {
        t.Errorf("Expect lexer.ch to be %q, but got %q", input[0], lexer.getCh())
    }
}

func TestNewLexerEmptyInput(t *testing.T) {
    var input = ""
    var lexer = NewLexer(input)

    if lexer.input != "" {
        t.Errorf("Expected lexer.input to be %q, but got %q", "", lexer.input)
    }
    if lexer.pos != 0 {
        t.Errorf("Expected lexer.position to be %q, but got %q", 0, lexer.pos)
    }
}

type ExpectedToken struct {
    Type    string
    Literal string
}

func checksForNextToken(lexer *Lexer, t *testing.T, expectedTokens []ExpectedToken) {
    for i, expected := range expectedTokens {
        var token = lexer.GetNextToken()

        if token.Type != expected.Type {
            t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
                i, expected.Type, token.Type)
        }

        if token.Literal != expected.Literal {
            t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
                i, expected.Literal, token.Literal)
        }
    }
}

func TestGetNextToken(t *testing.T) {
    var input = "=+(){},;"
    var lexer = NewLexer(input)

    var expectedTokens = []ExpectedToken {
        {token.Assign, "="},
        {token.Plus, "+"},
        {token.Lparen, "("},
        {token.Rparen, ")"},
        {token.Lbrace, "{"},
        {token.Rbrace, "}"},
        {token.Comma, ","},
        {token.Semicolon, ";"},
        {token.Eof, ""},
    }

    checksForNextToken(lexer, t, expectedTokens)
}

func TestGetNextToken2(t *testing.T) {
    var input = `let five = 5;`
    var lexer = NewLexer(input)

    var expectedTokens = []ExpectedToken {
        {token.Let, "let"},
        {token.Ident, "five"},
        {token.Assign, "="},
        {token.Int, "5"},
        {token.Semicolon, ";"},
    }

    checksForNextToken(lexer, t, expectedTokens)
}

func TestGetNextToken3(t *testing.T) {
    var input = `
        let five = 5;
        let ten = 10;

        let add = fn(x, y) {
            x + y;
        };

        let result = add(five, ten);
    `
    var lexer = NewLexer(input)

    var expectedTokens = []ExpectedToken {
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

    checksForNextToken(lexer, t, expectedTokens)
}

func TestGetNextToken4(t *testing.T) {
    var input = `
        !-/*5;
        5 < 10 > 5;
    `
    var lexer = NewLexer(input)

    var expectedTokens = []ExpectedToken {
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

    checksForNextToken(lexer, t, expectedTokens)
}

func TestGetNextToken5(t *testing.T) {
    var input = `
        if (5 < 10) {
            return true;
        } else {
            return false;
        }
    `
    var lexer = NewLexer(input)

    var expectedTokens = []ExpectedToken {
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

    checksForNextToken(lexer, t, expectedTokens)
}

func TestNextToken6(t *testing.T) {
    var input = `
        10 == 10;
        9 !=  10;
    `
    var lexer = NewLexer(input)

    var expectedTokens = []ExpectedToken {
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

    checksForNextToken(lexer, t, expectedTokens)
}

func TestNextTokenAllTokens(t *testing.T) {
    var input = `foo 5 = + - ! * / < > == != , ; ( ) { } fn let true false if else return`
    var lexer = NewLexer(input)

    var expectedTokens = []ExpectedToken {
        // Identifier + literals
        {token.Ident, "foo"},
        {token.Int, "5"},

        // Operators
        {token.Assign, "="},
        {token.Plus, "+"},
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

    checksForNextToken(lexer, t, expectedTokens)
}

func TestStrings(t *testing.T) {
    var input = `
        "foo bar";
        "foobar";
        "";
        "hello, world!";
    `
    var expectedTokens = []ExpectedToken {
        {token.String, "foo bar"},
        {token.Semicolon, ";"},
        {token.String, "foobar"},
        {token.Semicolon, ";"},
        {token.String, ""},
        {token.Semicolon, ";"},
        {token.String, "hello, world!"},
        {token.Semicolon, ";"},
    }

    var lexer = NewLexer(input)

    checksForNextToken(lexer, t, expectedTokens)
}

func TestArrays(t *testing.T) {
    var input = `[1, 2, 3]`
    var expectedTokens = []ExpectedToken {
        { token.Lbracket, "[" },
        { token.Int, "1" },
        { token.Comma, "," },
        { token.Int, "2" },
        { token.Comma, "," },
        { token.Int, "3" },
        { token.Rbracket, "]" },
    }
    var lexer = NewLexer(input)

    checksForNextToken(lexer, t, expectedTokens)
}

func TestHashs(t *testing.T) {
    var input = `{ "foo": "bar" }`
    var expectedTokens = []ExpectedToken {
        { token.Lbrace, "{"   },
        { token.String, "foo" },
        { token.Colon,  ":"   },
        { token.String, "bar" },
        { token.Rbrace, "}"   },
    }
    var lexer = NewLexer(input)

    checksForNextToken(lexer, t, expectedTokens)
}
