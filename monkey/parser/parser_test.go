// monkey/parser/parser_test.go

package parser

import (
    _"fmt"
    "bytes"
    "monkey/ast"
    "monkey/lexer"
    "monkey/token"
    "strconv"
    "testing"
)

func checkParserErrors(t *testing.T, par *Parser) {
    errors := par.Errors()
    if len(errors) == 0 { return }
    for i, err := range errors {
        t.Errorf("[%d] Parser Error: %s\n", i, err)
    }
    t.FailNow()
}

func TestLetStatement(t *testing.T) {
    input := `
        let x = 5;
        let y = 10;
        let z = 15;
    `
    lex := lexer.NewLexer(input)
    par := NewParser(lex)
    pro := par.ParseProgram()

    checkParserErrors(t, par)
    if pro == nil {
        t.Fatalf("Program is nil")
    }
    if len(pro.Statements) != 3 {
        t.Fatalf("Expected program to have %d statements but got %d\n", 3, len(pro.Statements))
    }

    tests := []struct { expectedIdentifier string } {
        {"x"}, {"y"}, {"z"},
    }

    for i, test := range tests {
        if pro.Statements[i].TokenLiteral() != "let" {
            t.Errorf("Expected statement literal to be '%s' but got '%s' instead",
                "let", pro.Statements[i].TokenLiteral())
        }
        stm, ok := pro.Statements[i].(*ast.LetStatement)
        if !ok {
            t.Errorf("Is not a LetStatement")
        }
        if stm.Token.Type != token.LET {
            t.Errorf("[%dl] Expected identifier to be %s but got %s", i, token.LET, stm.Token.Type)
        }
        if stm.Identifier.Value != test.expectedIdentifier {
            t.Errorf("[%d] Expected identifier to be %s but got %s", i, stm.Identifier.Value, test.expectedIdentifier)
        }
    }
}

func TestReturnStatement(t *testing.T) {
    input := `
        return 5;
        return 10;
        return 1234;
    `
    lex := lexer.NewLexer(input)
    par := NewParser(lex)
    pro := par.ParseProgram()

    checkParserErrors(t, par)
    if pro == nil {
        t.Fatalf("Program is nil")
    }
    if len(pro.Statements) != 3 {
        t.Fatalf("Expected program to have %d statements but got %d\n", 3, len(pro.Statements))
    }
    for i, currStm := range pro.Statements {
        stm, ok := currStm.(*ast.ReturnStatement)
        if !ok {
            t.Errorf("Is not a ReturnStatement")
        }
        if stm.Token.Type != token.RETURN {
            t.Errorf("[%d] Expected token type to be %s but got %s", i, token.RETURN, stm.Token.Type)
        }
        if stm.Token.Literal != "return" {
            t.Errorf("[%d] Expected token literal to be '%s' but got '%s'", i, "return", stm.Token.Literal)
        }
        if stm.TokenLiteral() != "return" {
            t.Errorf("[%d] Expected TokenLiteral to be '%s' but got '%s'", i, "return", stm.Token.Literal)
        }
    }
}

func TestErrorsOnLetStatement(t *testing.T) {
    input := `
        let x 5;
        let = 10;
        let 15;
    `
    lex := lexer.NewLexer(input)
    par := NewParser(lex)
    pro := par.ParseProgram()

    if pro == nil {
        t.Fatalf("Program is nil")
    }
    if len(par.errors) != 4 {
        t.Fatalf("Expected number of errors to be %d but got %d instead.", 4, len(par.errors))
    }
}

func TestIdentifierExpression(t *testing.T) {
    input := "foo;"
    lex := lexer.NewLexer(input)
    par := NewParser(lex)
    pro := par.ParseProgram()

    checkParserErrors(t, par)
    if len(pro.Statements) != 1 {
        t.Fatalf("Expected program statements length to be %d but got %d", 1, len(pro.Statements))
    }
    stm, ok := pro.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Not an expression statement")
    }
    ident, ok := stm.Expression.(*ast.Identifier)
    if !ok {
        t.Fatalf("Statement expression is not an identifier")
    }
    if ident.Value != "foo" {
        t.Errorf("Expected identifier value to be '%s' but got '%s' instead", "foo", ident.Value)
    }
    if ident.TokenLiteral() != "foo" {
        t.Errorf("Expected statement expression to be '%s' but got '%s' instead", "foo", ident.TokenLiteral())
    }
}

func TestIntegerLiteralExpression(t *testing.T) {
    input := "5;"
    lex := lexer.NewLexer(input)
    par := NewParser(lex)
    pro := par.ParseProgram()

    if len(pro.Statements) != 1 {
        t.Fatalf("Expected program number of statements to be %d but got %d instead", 1, len(pro.Statements))
    }
    stm, ok := pro.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Expected an ast.ExpressionStatement but got %T instead", pro.Statements[0])
    }
    intLit, ok := stm.Expression.(*ast.IntegerLiteral)
    if !ok {
        t.Fatalf("Expected an ast.IntegerLiteral but got %T instead", stm.Expression)
    }
    if intLit.TokenLiteral() != "5" {
        t.Errorf("Expected statement expression to be '%s' but got '%s' instead", "5", intLit.TokenLiteral())
    }
    if intLit.Token.Type != token.INT {
        t.Errorf("Expected integer literal token type to be %s but got %s instead", token.INT, intLit.Token.Type)
    }
    if intLit.Value != 5 {
        t.Errorf("Expected integer literal value to be %d but got %d instead", 5, intLit.Value)
    }
}

func TestParsingPrefixExpression(t *testing.T) {
    input := `
        !5;
        -15;
    `
    lex := lexer.NewLexer(input)
    par := NewParser(lex)
    pro := par.ParseProgram()

    if pro == nil {
        t.Fatalf("Program is nil")
    }
    if len(pro.Statements) != 2 {
        t.Fatalf("Expected program number of statements to be %d but got %d instead", 2, len(pro.Statements))
    }
    checkParserErrors(t, par)

    tests := []struct { operator string; integerValue int64 } {
        {"!", 5}, {"-", 15},
    }
    for i, test := range tests {
        stm, ok := pro.Statements[i].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("Not an expression statement")
        }
        exp, ok := stm.Expression.(*ast.PrefixExpression)
        if !ok {
            t.Fatalf("Expression statement is not a prefix statement")
        }
        if exp.Operator != test.operator {
            t.Errorf("Expected operator to be '%s' but got '%s' instead", test.operator, exp.Operator)
        }
        val, err := strconv.ParseInt(exp.Right.String(), 10, 64)
        if err != nil {
            t.Fatalf("Cound not convert to int64")
        }
        if val != test.integerValue {
            t.Errorf("Expected prefix expression value to be %d but got %d instead", test.integerValue, val)
        }
    }
}

func TestParsingInfixExpression(t *testing.T) {
    tests := []struct {
        input string
        leftValue int64
        operator string
        rightValue int64
    } {
        { "5 + 5",  5, "+",  5 },
        { "5 - 5",  5, "-",  5 },
        { "5 * 5",  5, "*",  5 },
        { "5 / 5",  5, "/",  5 },
        { "5 > 5",  5, ">",  5 },
        { "5 < 5",  5, "<",  5 },
        { "5 == 5", 5, "==", 5 },
        { "5 != 5", 5, "!=", 5 },
    }
    var acc bytes.Buffer
    for _, test := range tests {
        acc.WriteString(test.input + ";\n")
    }
    input := acc.String()
    lex := lexer.NewLexer(input)
    par := NewParser(lex)
    pro := par.ParseProgram()

    checkParserErrors(t, par)
    if pro == nil {
        t.Fatalf("Program is nil")
    }
    if len(pro.Statements) != 8 {
        t.Fatalf("Expected program number of statements to be %d but got %d instead", 8, len(pro.Statements))
    }

    for i, test := range tests {
        stm, ok := pro.Statements[i].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("[%d] Not an expression statement", i)
        }
        expression, ok := stm.Expression.(*ast.InfixExpression)
        if !ok {
            t.Fatalf("[%d] Expression statement is not a infix expression", i)
        }
        // Test Left
        testIntegerLiteral(t, expression.Left, test.leftValue)
        // Operator
        if expression.Operator != test.operator {
            t.Fatalf("[%d] Expected Infix Expression Operator to be '%s' but got '%s' instead",
                i, test.operator, expression.Operator)
        }
        // Test Right
        testIntegerLiteral(t, expression.Right, test.rightValue)
    }
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, expectedValue int64) {
    intLiteral, ok := exp.(*ast.IntegerLiteral)
    if !ok {
        t.Errorf("Expression is not an Integer Literal")
        return
    }
    if intLiteral.Value != expectedValue {
        t.Errorf("Expected Integer Literal value to be %d but got %d instead", expectedValue, intLiteral.Value)
    }
    expectedStr := strconv.FormatInt(expectedValue, 10)
    if intLiteral.TokenLiteral() != expectedStr {
        t.Errorf("Expected Integer token literal to be '%s' but got '%s' instead", expectedStr, intLiteral.TokenLiteral())
    }
}

func TestOperatorPrecedenceParsing(t *testing.T) {
    tests := []struct {
          input string;                 expected string
    } {
        { "-a * b",                     "((-a) * b)" },
        { "!-a",                        "(!(-a))" },
        { "a + b + c",                  "((a + b) + c)" },
        { "a + b - c",                  "((a + b) - c)" },
        { "a * b * c",                  "((a * b) * c)" },
        { "a * b / c",                  "((a * b) / c)" },
        { "a + b / c",                  "(a + (b / c))" },
        { "a + b * c + d / e - f",      "(((a + (b * c)) + (d / e)) - f)" },
        { "3 + 4; -5 * 5",              "(3 + 4)((-5) * 5)" },
        { "5 > 4 == 3 < 4",             "((5 > 4) == (3 < 4))" },
        { "5 < 4 != 3 > 4",             "((5 < 4) != (3 > 4))" },
        { "3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))" },
    }

    for _, test := range tests {
        lex := lexer.NewLexer(test.input)
        par := NewParser(lex)
        pro := par.ParseProgram()

        checkParserErrors(t, par)

        progStr := pro.String()
        if progStr != test.expected {
            t.Errorf("Expected program string to be '%s' but got '%s' instead", test.expected, progStr)
        }
    }
}
