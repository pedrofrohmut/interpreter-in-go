// monkey/parser/parser_test.go

package parser

import (
    "testing"
    "monkey/lexer"
)

func TestLetStatements(t *testing.T) {
    input := `
        let x = 5;
        let y = 10;
        let z = 15;
    `
    lexer := lexer.NewLexer(input)
    parser := NewParser(lexer)
    program := parser.ParseProgram()

    if len(program.Statements) != 3 {
        t.Fatalf("Expected program to have %d statements but got %d instead", 3, len(program.Statements))
    }

    program.PrintStatements()
}
