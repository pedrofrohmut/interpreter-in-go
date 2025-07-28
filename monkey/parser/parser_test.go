// monkey/parser/parser_test.go

package parser

import (
    "testing"
    "reflect"
    "monkey/ast"
    "monkey/lexer"
    "monkey/token"
)

func TestLetStatement(t *testing.T) {
    input := `
        let x = 5;
        let y = 10;
        let z = 15;
    `
    lex := lexer.NewLexer(input)
    par := NewParser(lex)
    pro := par.ParseProgram()

    if pro == nil {
        t.Fatalf("Program is nill")
    }
    if len(pro.Statements) != 3 {
        t.Fatalf("Expected program to have %d statements but got %d\n", 3, len(pro.Statements))
    }

    tests := []struct { expectedIdentifier string } {
        {"x"}, {"y"}, {"z"},
    }

    for i, test := range tests {
        stm, ok := pro.Statements[i].(*ast.LetStatement)
        if !ok {
            t.Errorf("Is not a LetStatement")
        }
        if stm.Token.Type != token.LET {
            t.Errorf("[%dl] Expected identifier to be %s but got %s", i, token.LET, stm.Token.Type)
        }
        if stm.Identifier != test.expectedIdentifier {
            t.Errorf("[%d] Expected identifier to be %s but got %s", i, stm.Identifier, test.expectedIdentifier)
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

    if pro == nil {
        t.Fatalf("Program is nill")
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
        t.Fatalf("Program is nill")
    }

    for i, tmp := range pro.Statements {
        if tmp != nil && !reflect.ValueOf(tmp).IsNil() {
            t.Errorf("[%d] Current statement is not nil as expected for an invalid statement", i)
        }
    }

    expectedErrCount := 3
    if len(par.errors) < expectedErrCount {
        t.Fatalf("Expected number of errors to be %d but got %d instead.", expectedErrCount, len(par.errors))
    }
}
