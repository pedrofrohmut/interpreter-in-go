// monkey/ast/ast.go

package ast

import (
    "monkey/token"
)

type Node interface {
    TokenLiteral() string
}

type Statement interface {
    Node
    statementNode()
}

type Expression interface {
    Node
    expressionNode()
}

type Program struct {
    Statements []Statement
}

func NewProgram() *Program {
    return &Program {
        Statements: []Statement{},
    }
}

func (pr * Program) TokenLiteral() string {
    if len(pr.Statements) == 0 {
        return ""
    }
    return pr.Statements[0].TokenLiteral()
}

type Identifier struct {
    Token token.Token
    Value string
}

func (iden *Identifier) expressionNode() {}

func (iden *Identifier) TokenLiteral() string {
    return iden.Token.Literal
}

type LetStatement struct {
    Token token.Token
    Name *Identifier
    Value Expression
}

func NewLetStatement(tk token.Token) *LetStatement {
    return &LetStatement { Token: tk }
}

func (stm *LetStatement) statementNode() {}

func (stm *LetStatement) TokenLiteral() string {
    return stm.Token.Literal
}

type ReturnStatement struct {
    Token token.Token
    ReturnValue Expression
}

func NewReturnStatement(tk token.Token) *ReturnStatement {
    return &ReturnStatement { Token: tk }
}

func (stm *ReturnStatement) statementNode() {}

func (stm *ReturnStatement) TokenLiteral() string {
    return stm.Token.Literal
}
