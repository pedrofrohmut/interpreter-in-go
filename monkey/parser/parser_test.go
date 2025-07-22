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

func TestLetStatements(t *testing.T) {
    input := `
        let x = 5;
        let y = 10;
        let foobar = 838383;
    `
    lx := lexer.NewLexer(input)
    par := NewParser(lx)
    program := par.ParseProgram()

    if program == nil {
        t.Fatalf("ParseProgram() returned nil")
    }
    if len(program.Statements) != 3 {
        t.Fatalf("program.Statements expected 3 but got=%d", len(program.Statements))
    }

    tests := []ExpectedIdentifier {
        {"x"},
        {"y"},
        {"foobar"},
    }

    checkLetStatements(program, tests, t)
}
