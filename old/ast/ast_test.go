// monkey/ast/ast_test.go

package ast

import (
    "monkey/token"
    "testing"
)

func TestString(t *testing.T) {
    expectedProgram := "let myVar = anotherVar;"
    pro := NewProgram()
    // Hard coded AST example
    pro.Statements = []Statement {
        &LetStatement {
            Token: token.NewTokenStr(token.LET, "let"),
            Identifier: NewIdentifier(token.NewTokenStr(token.IDENT, "myVar"), "myVar"),
            Expression: NewIdentifier(token.NewTokenStr(token.IDENT, "anotherVar"), "anotherVar"),
        },
    }
    programString := pro.String()

    if programString != expectedProgram {
        t.Errorf("Expected program to string to be '%s' but got '%s' instead", expectedProgram, programString)
    }
}
