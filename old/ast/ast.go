// monkey/ast/ast.go

package ast

import (
    "bytes"
    "monkey/token"
)

type Node interface {
    TokenLiteral() string
    // Utility method
    String() string
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

func (pr *Program) TokenLiteral() string {
    if len(pr.Statements) == 0 {
        return ""
    }
    return pr.Statements[0].TokenLiteral()
}

func (pr *Program) String() string {
    var out bytes.Buffer
    for _, s := range pr.Statements {
        out.WriteString(s.String())
    }
    return out.String()
}

type Identifier struct {
    Token token.Token
    Value string
}

func NewIdentifier(tk token.Token) *Identifier {
    return &Identifier { Token: tk, Value: tk.Literal }
}

func (iden *Identifier) String() string {
    return iden.Value
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

func (stm *LetStatement) String() string {
    var out bytes.Buffer
    out.WriteString(stm.TokenLiteral() + " ")
    out.WriteString(stm.Name.String())
    out.WriteString(" = ")
    if stm.Value != nil {
        out.WriteString(stm.Value.String())
    }
    out.WriteString(";")
    return out.String()
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

func (stm *ReturnStatement) String() string {
    var out bytes.Buffer
    out.WriteString(stm.TokenLiteral() + " ")
    if stm.ReturnValue != nil {
        out.WriteString(stm.ReturnValue.String())
    }
    out.WriteString(";")
    return out.String()
}

func (stm *ReturnStatement) statementNode() {}

func (stm *ReturnStatement) TokenLiteral() string {
    return stm.Token.Literal
}

type ExpressionStatement struct {
    Token token.Token
    Expression Expression
}

func NewExpressionStatement(tk token.Token) *ExpressionStatement {
    return &ExpressionStatement { Token: tk }
}

func (stm *ExpressionStatement) String() string {
    if stm.Expression != nil {
        return stm.Expression.String()
    }
    return ""
}

func (ex *ExpressionStatement) statementNode() {}

func (ex *ExpressionStatement) TokenLiteral() string {
    return ex.TokenLiteral()
}
