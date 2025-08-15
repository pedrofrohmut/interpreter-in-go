// monkey/ast/ast.go

package ast

import (
    "fmt"
    "strconv"
)

type Node interface {
    node()
    String() string
}

type Statement interface {
    Node
    statement()
}

type Expression interface {
    Node
    expression()
}

type Program struct {
    Statements []Statement
}

// @Impl
func (this *Program) node() {}

func NewProgram() *Program {
    return &Program {
        Statements: []Statement {},
    }
}

func (this *Program) PrintStatements() {
    // TODO: Discover why this bullshit is not working
    if len(this.Statements) == 0 { return }
    for i, stm := range this.Statements {
        if stm != nil {
            fmt.Printf("[%d] %s\n", i, stm.String())
        } else {
            fmt.Printf("[%d] nil_statement\n", i)
        }
    }
}

type LetStatement struct {
    Identifier string
    Expression Expression
}

// @Impl
func (this *LetStatement) node() {}

// @Impl
func (this *LetStatement) statement() {}

// @Impl
func (this *LetStatement) String() string {
    expression := "%TODO%"
    return "let " + this.Identifier + " = " + expression + ";"
}

func NewLetStatement() *LetStatement {
    return &LetStatement {}
}

type ReturnStatement struct {
    Expression Expression
}

// @Impl
func (this *ReturnStatement) node() {}

// @Impl
func (this *ReturnStatement) statement() {}

// @Impl
func (this *ReturnStatement) String() string {
    expression := "%TODO%"
    if expression == "" {
        return "return;"
    }
    return "return " + expression + ";"
}

func NewReturnStatement() *ReturnStatement {
    return &ReturnStatement {}
}

type ExpressionStatement struct {
    Expression Expression
}

// @Impl
func (this *ExpressionStatement) node() {}

// @Impl
func (this *ExpressionStatement) statement() {}

// @Impl
func (this *ExpressionStatement) String() string {
    return this.Expression.String() + ";"
}

func NewExpressionStatement() *ExpressionStatement {
    return &ExpressionStatement {}
}

type Identifier struct {
    Value string
}

// @Impl
func (this *Identifier) node() {}

// @Impl
func (this *Identifier) expression() {}

// @Impl
func (this *Identifier) String() string { return this.Value }

func NewIdentifier(value string) *Identifier {
    return &Identifier { Value: value }
}

type IntegerLiteral struct {
    Value int64
}

// @Impl
func (this *IntegerLiteral) node() {}

// @Impl
func (this *IntegerLiteral) expression() {}

// @Impl
func (this *IntegerLiteral) String() string { return strconv.FormatInt(this.Value, 10) }

func NewIntegerLiteral(value int64) *IntegerLiteral {
    return &IntegerLiteral { Value: value }
}

type PrefixExpression struct {
    Operator string
    Value Expression
}

// @Impl
func (this *PrefixExpression) node() {}

// @Impl
func (this *PrefixExpression) expression() {}

// @Impl
func (this *PrefixExpression) String() string {
    return "(" + this.Operator + this.Value.String() + ")"
}

func NewPrefixExpression(operator string, value Expression) *PrefixExpression {
    return &PrefixExpression { Operator: operator, Value: value }
}

type InfixExpression struct {
    Operator string
    Left Expression
    Right Expression
}

// @Impl
func (this *InfixExpression) node() {}

// @Impl
func (this *InfixExpression) expression() {}

// @Impl
func (this *InfixExpression) String() string {
    return "(" + this.Left.String() + " " + this.Operator + " " + this.Right.String() + ")"
}

func NewInfixExpression(left Expression) *InfixExpression {
    return &InfixExpression { Left: left }
}

type Boolean struct {
    Value bool
}

// @Impl
func (this *Boolean) node() {}

// @Impl
func (this *Boolean) expression() {}

// @Impl
func (this *Boolean) String() string {
    return strconv.FormatBool(this.Value)
}
