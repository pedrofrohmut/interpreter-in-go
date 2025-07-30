// monkey/ast/ast.go
/*
     AST - Abstract Syntax Tree. Is the data structure the results for parsing
   the source code of monkey language
*/

package ast

import (
    "bytes"
    "monkey/token"
)

type Node interface {
    String() string
}

type Statement interface {
    Node
}

type Program struct {
    Node
    Statements []Statement
}

// @Impl
func (pro *Program) String() string {
    var out bytes.Buffer
    for _, stm := range pro.Statements {
        out.WriteString(stm.String())
    }
    return out.String()
}

func NewProgram() *Program {
    return &Program { Statements: [] Statement{} }
}

// TODO: Keeping checking if Identifier needs to be a separated struct for reuse
type LetStatement struct {
    Statement
    Token token.Token
    Identifier string
    Expression string
}

// @Impl
func (stm *LetStatement) String() string {
    var out bytes.Buffer
    out.WriteString(stm.Token.Literal + " ")
    out.WriteString(stm.Identifier + " = ")
    if stm.Expression != "" {
        out.WriteString(stm.Expression)
    }
    out.WriteString(";")
    return out.String()
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

// @Impl
func (stm *ReturnStatement) String() string {
    var out bytes.Buffer
    out.WriteString(stm.Token.Literal + " ")
    if stm.Expression != "" {
        out.WriteString(stm.Expression)
    }
    out.WriteString(";")
    return out.String()
}

func NewReturnStatement() *ReturnStatement {
    return &ReturnStatement {
        Token: token.NewTokenStr(token.RETURN, "return"),
    }
}

// We are having Expression Statements because in monkey you can have expressions
// as statements. Exp: 5 * 5 + 3;. So it is needed to have it as statement here
type ExpressionStatement struct {
    Statement
    Token token.Token // The first token of the expression
    Expression string
}

// @Impl
func (stm *ExpressionStatement) String() string {
    var out bytes.Buffer
    if stm.Expression != "" {
        out.WriteString(stm.Expression)
    }
    // TODO: Check if should just return an empty string or ; is fine
    out.WriteString(";")
    return out.String()
}

func NewExpressionStatement(first token.Token) *ExpressionStatement {
    return &ExpressionStatement { Token: first }
}

type IntegerLiteral struct {
    Node
    Token token.Token
    Value int64
}

func NewIntegerLiteral(val int64) *IntegerLiteral {
    return &IntegerLiteral {
        Token: token.Token { Type: token.INT, Literal: string(val) },
        Value: val,
    }
}

// @Impl
func (lit *IntegerLiteral) String() string {
    return lit.Token.Literal
}
