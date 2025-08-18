// monkey/ast/ast.go

package ast

import (
    "fmt"
    "bytes"
    "strconv"
    "strings"
    "monkey/utils"
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
    for i, stm := range this.Statements {
        fmt.Printf("[%d] %s;\n", i, stm.String())
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
    var out bytes.Buffer
    out.WriteString("let ")
    out.WriteString(this.Identifier)
    out.WriteString(" = ")
    out.WriteString(this.Expression.String())
    return out.String()
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
    if utils.IsNill(this.Expression) { return "return" }

    var out bytes.Buffer
    out.WriteString("return ")
    out.WriteString(this.Expression.String())
    return out.String()
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
    return this.Expression.String()
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

type StatementsBlock struct {
    Statements []Statement
}

// @Impl
func (this *StatementsBlock) node() {}

// @Impl
func (this *StatementsBlock) statement() {}

// @Impl
func (this *StatementsBlock) String() string {
    var out bytes.Buffer
    for _, stm := range this.Statements {
        out.WriteString(stm.String())
    }
    return out.String()
}

type IfExpression struct {
    Condition Expression
    ConsequenceBlock *StatementsBlock
    AlternativeBlock *StatementsBlock
}

// @Impl
func (this *IfExpression) node() {}

// @Impl
func (this *IfExpression) expression() {}

// @Impl
func (this *IfExpression) String() string {
    var out bytes.Buffer
    out.WriteString("if " + this.Condition.String() + " {")

    out.WriteString(this.ConsequenceBlock.String())

    if !utils.IsNill(this.AlternativeBlock) {
        out.WriteString("} else {")
        out.WriteString(this.AlternativeBlock.String())
    }

    out.WriteString("}")
    return out.String()
}

type FunctionLiteral struct {
    Arguments []Identifier
    Body *StatementsBlock
}

// @Impl
func (this *FunctionLiteral) node() {}

// @Impl
func (this *FunctionLiteral) expression() {}

// @Impl
func (this *FunctionLiteral) String() string {
    var out bytes.Buffer

    var args = []string {}
    for _, arg := range this.Arguments {
        args = append(args, arg.String())
    }

    var stms = []string {}
    for _, stm := range this.Body.Statements {
        stms = append(stms, stm.String())
    }

    out.WriteString("fn (")
    if len(args) > 0 {
        out.WriteString(strings.Join(args, ", "))
    }
    out.WriteString(") {")
    if len(stms) > 0 {
        out.WriteString(" " + strings.Join(stms, " ") + " ")
    }
    out.WriteString("}")

    return out.String()
}

type CallExpression struct {
    Expression Expression
    Parameters []Expression
}

// @Impl
func (this *CallExpression) node() {}

// @Impl
func (this *CallExpression) expression() {}

// @Impl
func (this *CallExpression) String() string {
    var out bytes.Buffer

    var params = []string {}
    for _, param := range this.Parameters {
        params = append(params, param.String())
    }

    out.WriteString(this.Expression.String() + "(")
    out.WriteString(strings.Join(params, ", "))
    out.WriteString(")")

    return out.String()
}
