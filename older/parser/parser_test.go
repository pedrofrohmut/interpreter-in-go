// monkey/parser/parser_test.go

package parser

import (
    "testing"
    "monkey/lexer"
    "monkey/ast"
)

type ExpectedIdentifier struct {
    expectedIdentifier string
}

func checkLetStatements(program *ast.Program, tests []ExpectedIdentifier, t *testing.T) {
    for i, test := range tests {
        stm := program.Statements[i]
        if stm.TokenLiteral() != "let" {
            t.Fatalf("Expected the first token to be 'let' but got=%s\n", stm.TokenLiteral())
        }

        letStm, ok := stm.(*ast.LetStatement)
        if !ok {
            t.Fatalf("Current statement is not a LetStatement. Expected ast.LetStatement but got=%s\n", stm)
        }

        if letStm.Name.Value != test.expectedIdentifier {
            t.Fatalf("Expected %s but got %s\n", test.expectedIdentifier, letStm.Name.Value)
        }

        if letStm.Name.TokenLiteral() != test.expectedIdentifier {
            t.Fatalf("Expected statement name to be %s but got %s\n", test.expectedIdentifier, letStm.Name.TokenLiteral())
        }
    }
}

func checkParserErrors(t *testing.T, par *Parser) {
    errors := par.GetErrors()
    if len(errors) == 0 { return }
    t.Errorf("Parser has %d errors.\n", len(errors))
    for _, msg := range errors {
        t.Errorf("Parser error: %q\n", msg)
    }
    t.FailNow()
}

func checkProgram(program *ast.Program, t *testing.T) {
    if program == nil {
        t.Fatalf("ParseProgram() returned nil")
    }
    if len(program.Statements) != 3 {
        t.Fatalf("program.Statements expected 3 but got=%d", len(program.Statements))
    }
}

func TestLetStatements(t *testing.T) {
    input := `
        let x = 5;
        let y = 10;
        let foobar = 838383;
    `
    lx := lexer.NewLexer(input)
    par := NewParser(lx)
    program := par.ParseProgram()
    checkParserErrors(t, par)
    checkProgram(program, t)

    tests := []ExpectedIdentifier {
        {"x"},
        {"y"},
        {"foobar"},
    }

    checkLetStatements(program, tests, t)
}

func TestReturnStatement(t *testing.T) {
    input := `
        return 5;
        return 10;
        return 12345;
    `
    lx := lexer.NewLexer(input)
    par := NewParser(lx)
    program := par.ParseProgram()
    checkParserErrors(t, par)
    checkProgram(program, t)

    for _, currStm := range program.Statements {
        stm, ok := currStm.(*ast.ReturnStatement)
        if !ok {
            t.Errorf("Expected a return statement but got %T", stm)
            continue
        }
        if stm.TokenLiteral() != "return" {
            t.Errorf("Expected token literal to be 'return' but got %q", stm.TokenLiteral())
        }
    }
}

func TestIdentifierExpression(t *testing.T) {
    input := `foobar;`
    lx := lexer.NewLexer(input)
    par := NewParser(lx)
    program := par.ParseProgram()
    checkParserErrors(t, par)

    if len(program.Statements) != 1 {
        t.Fatalf("Program does not have right number of statements, got %d", len(program.Statements))
    }
    stm, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("The first statement is not an ast.ExpressionStatement, got %T", program.Statements[0])
    }

    iden, ok := stm.Expression.(*ast.Identifier)
    if !ok {
        t.Fatalf("It is not an identifier, got %T", stm.Expression)
    }
    if iden.Value != "foobar" {
        t.Errorf("The identifier is not %s, got %s", "foobar", iden.Value)
    }
    if iden.TokenLiteral() != "foobar" {
        t.Errorf("Expected identifier literal to be %s but got %s", "foobar", iden.TokenLiteral())
    }
}
