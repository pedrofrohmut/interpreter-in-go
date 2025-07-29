// monkey/ast/ast_test.go

package ast

import (
    "monkey/token"
    "testing"
)

func TestString(t *testing.T) {
    pro := NewProgram()
    // Hard coded AST example
    pro.Statements = []Statement {
        &LetStatement {
            Token: token.NewTokenStr(token.LET, "let"),
            Identifier: "myVar",
            Expression: "anotherVar",
        },
    }
    programString := pro.String()

    expectedProgram := "let myVar = anotherVar;"
    if programString != expectedProgram {
        t.Errorf("Expected program to string to be '%s' but got '%s' instead", expectedProgram, programString)
    }
}
