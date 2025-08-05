// monkey/ast/ast.go
/*
     AST - Abstract Syntax Tree. Is the data structure the results for parsing
   the source code of monkey language
*/

package ast

import (
    "bytes"
    _"strconv"
    "monkey/token"
)

type Node interface {
    TokenLiteral() string
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

// @Impl
func (this *Program) TokenLiteral() string {
    if len(this.Statements) == 0 { return "" }
    return this.Statements[0].TokenLiteral()
}

// @Impl
func (this *Program) String() string {
    var out bytes.Buffer
    for _, stm := range this.Statements {
        out.WriteString(stm.String())
    }
    return out.String()
}

func NewProgram() *Program {
    return &Program {
        Statements: [] Statement{},
    }
}

type Identifier struct {
    Token token.Token
    Value string
}

// @Impl
func (this *Identifier) expressionNode() {}

// @Impl
func (this *Identifier) TokenLiteral() string { return this.Token.Literal }

// @Impl
func (this *Identifier) String() string { return this.Value }

func NewIdentifier(tok token.Token, val string) *Identifier {
    return &Identifier { Token: tok, Value: val }
}

type LetStatement struct {
    Token token.Token
    Identifier *Identifier
    Expression Expression
}

// @Impl
func (this *LetStatement) statementNode() {}

// @Impl
func (this *LetStatement) TokenLiteral() string { return this.Token.Literal }

// @Impl
func (this *LetStatement) String() string {
    var out bytes.Buffer
    out.WriteString(this.Token.Literal + " ")
    out.WriteString(this.Identifier.Value + " = ")
    if this.Expression != nil {
        out.WriteString(this.Expression.String())
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
    Token token.Token
    Expression Expression
}

// @Impl
func (this *ReturnStatement) statementNode() {}

// @Impl
func (this *ReturnStatement) TokenLiteral() string { return this.Token.Literal }

// @Impl
func (this *ReturnStatement) String() string {
    var out bytes.Buffer
    out.WriteString(this.Token.Literal + " ")
    if this.Expression != nil {
        out.WriteString(this.Expression.String())
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
    Token token.Token // The first token of the expression
    Expression Expression
}

// @Impl
func (this *ExpressionStatement) statementNode() {}

// @Impl
func (this *ExpressionStatement) TokenLiteral() string { return this.Token.Literal }

// @Impl
func (this *ExpressionStatement) String() string {
    if this.Expression == nil {
        return ""
    }
    return this.Expression.String()
}

func NewExpressionStatement(tok token.Token) *ExpressionStatement {
    return &ExpressionStatement { Token: tok }
}

type IntegerLiteral struct {
    Token token.Token
    Value int64
}

// @Impl
func (this *IntegerLiteral) expressionNode() {}

// @Impl
func (this *IntegerLiteral) TokenLiteral() string { return this.Token.Literal }

// @Impl
func (this *IntegerLiteral) String() string { return this.Token.Literal }

func NewIntegerLiteral(tok token.Token, val int64) *IntegerLiteral {
    return &IntegerLiteral { Token: tok, Value: val }
}

type Boolean struct {
    Token token.Token
    Value bool
}

// @Impl
func (this *Boolean) expressionNode() {}

// @Impl
func (this *Boolean) TokenLiteral() string { return this.Token.Literal }

// @Impl
func (this *Boolean) String() string { return this.Token.Literal }

func NewBoolean(tok token.Token, val bool) *Boolean {
    return &Boolean { Token: tok, Value: val }
}

type PrefixExpression struct {
    Token token.Token
    Operator string
    Right Expression
}

// @Impl
func (this *PrefixExpression) expressionNode() {}

// @Impl
func (this *PrefixExpression) TokenLiteral() string { return this.Token.Literal }

// @Impl
func (this *PrefixExpression) String() string {
    var out bytes.Buffer
    out.WriteString("(")
    out.WriteString(this.Operator)
    out.WriteString(this.Right.String())
    out.WriteString(")")
    return out.String()
}

func NewPrefixExpression(token token.Token, operator string) *PrefixExpression {
    return &PrefixExpression { Token: token, Operator: operator }
}

type InfixExpression struct {
    Token token.Token // Operator token
    Left Expression
    Operator string
    Right Expression
}

// @Impl
func (this *InfixExpression) expressionNode() {}

// @Impl
func (this *InfixExpression) TokenLiteral() string { return this.Token.Literal }

// @Impl
func (this *InfixExpression) String() string {
    var out bytes.Buffer
    out.WriteString("(")
    out.WriteString(this.Left.String())
    out.WriteString(" " + this.Operator + " ")
    out.WriteString(this.Right.String())
    out.WriteString(")")
    return out.String()
}

func NewInfixExpression(tok token.Token, left Expression) *InfixExpression {
    return &InfixExpression { Token: tok, Operator: tok.Literal, Left: left }
}
