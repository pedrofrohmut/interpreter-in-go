// monkey/ast/ast.go

package ast

import (
    "fmt"
    _"bytes"
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
        fmt.Printf("[%d] %s\n", i, stm)
    }
}

type LetStatement struct {
    Identifier string
    Expression *Expression
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
