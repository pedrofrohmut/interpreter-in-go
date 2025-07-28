// monkey/parser/parser_test.go

package parser

import (
    "testing"
    // "fmt"
    "monkey/lexer"
    "monkey/token"
    "monkey/ast"
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
    if par.curr.Type != token.EOF {
        t.Fatalf("Parse did not get to the of the file.")
    }

    tests := []struct { expectedIdentifier string } {
        {"x"}, {"y"}, {"z"},
    }

    for i, test := range tests {
        stm, ok := pro.Statements[i].(*ast.LetStatement)
        if !ok {
            t.Fatalf("Is not a LetStatement")
        }

        if stm.Identifier.Value != test.expectedIdentifier {
            t.Fatalf("[%d] Expected identifier to be %s but got %s",
                i, stm.Identifier.Value, test.expectedIdentifier)
        }

        // fmt.Println(stm.GetTokenLiteral())
    }
}

func TestReturnStatements(t *testing.T) {
    input := `
        return 5;
        return 10;
        return 1 + 5;
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
    if par.curr.Type != token.EOF {
        t.Fatalf("Parse did not get to the of the file.")
    }

    // tests := []struct { expectedValue string } {
    //     {"5"}, {"10"}, {"1 + 5"},
    // }

    // for i, _ := range tests {
    //     fmt.Println(i)
    //     fmt.Println(pro.Statements[i])
    //
    //     _, ok := pro.Statements[i].(*ast.ReturnStatement)
    //     if !ok {
    //         t.Fatalf("Is not a Return Statement")
    //     }
    //     // if stm == nil {
    //     //     t.Fatalf("This statement is nil")
    //     // }
    //     //
    //     // if stm.Token.Type != token.RETURN {
    //     //     t.Fatalf("Expected token type to be %s but got %s", token.RETURN, stm.Token.Type)
    //     // }
    //     // if stm.Expression.GetTokenLiteral() != test.expectedValue {
    //     //     t.Fatalf("Expected expression to be '%s' but got '%s'", test.expectedValue, stm.Expression)
    //     // }
    // }
}
