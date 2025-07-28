// monkey/ast/ast.go
/*
      AST - Abstract Syntax Tree. Is the data structure the results for parsing
    the source code of monkey language
*/

package ast

import (
    "monkey/token"
)

// Check if this interface is necessary and also the GetTokenLiteral
type Node interface {
    // GetTokenLiteral() string
    String() string
}

type Expression interface {
    Node
}

type Statement interface {
    Node
}

// AST Root
type Program struct {
    Statements []Statement
}

func NewProgram() *Program {
    return &Program { Statements: []Statement {} }
}

// TODO: check if this method is necessary
// Also check why it is first statement only
// func (pr *Program) GetTokenLiteral() string {
//     if len(pr.Statements) < 1 {
//         return ""
//     }
//     return pr.Statements[0].GetTokenLiteral()
// }

// TODO: Maybe it doesnt need to be token. Only TokenType
type Identifier struct {
    Token token.Token // the token.IDENT token
    Value string
}

func NewIdentifier(value string) *Identifier {
    return &Identifier {
        Token: token.NewTokenStr(token.IDENT, value),
        Value: value,
    }
}

// TODO: Maybe you dont need a separated struct for identifier
// TODO: Maybe it doesnt need to be token. Only TokenType
type LetStatement struct {
    Statement
    Token token.Token // the token.LET token
    Identifier *Identifier
    Expression Expression
}

func NewLetStatement() *LetStatement {
    return &LetStatement {
        Token: token.NewTokenStr(token.LET, "let"),
    }
}

type ReturnStatement struct {
    Statement
    Token token.Token
    Expression Expression
}

func NewReturnStatement() *ReturnStatement {
    return &ReturnStatement {
        Token: token.NewTokenStr(token.RETURN, "return"),
    }
}
