// monkey/ast/ast.go
/*
      AST - Abstract Syntax Tree. Is the data structure the results for parsing
    the source code of monkey language
*/

package ast

import (
    "monkey/token"
)

type Statement any

type Program struct {
    Statements []Statement
}

func NewProgram() *Program {
    return &Program { Statements: [] Statement{} }
}

type LetStatement struct {
    Statement
    Token token.Token
    Identifier string
    Expression string
}

func NewLetStatement() *LetStatement {
    return &LetStatement {
        Token: token.NewTokenStr(token.LET, "let"),
    }
}

type ReturnStatement struct {
    Statement
    Token token.Token
    Expression string
}

func NewReturnStatement() *ReturnStatement {
    return &ReturnStatement {
        Token: token.NewTokenStr(token.RETURN, "return"),
    }
}
