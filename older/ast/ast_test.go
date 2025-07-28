// monkey/ast/ast_test.go

package ast

import (
    "testing"
    "monkey/token"
)

func TestString(t *testing.T) {
    program := &Program {
        Statements: []Statement {
            &LetStatement {
                Token: token.Token { Type: token.LET, Literal: "let" },
                Name: &Identifier {
                    Token: token.Token { Type: token.IDENT, Literal: "myVar" },
                    Value: "myVar",
                },
                Value: &Identifier {
                    Token: token.Token { Type: token.IDENT, Literal: "anotherVar" },
                    Value: "anotherVar",
                },
            },
        },
    }

    expectedOutput := "let myVar = anotherVar;"
    if program.String() != expectedOutput {
        t.Errorf("Expected program to string to be %s but got %s", expectedOutput, program.String())
    }
}
