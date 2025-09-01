// monkey/object/object.go

package object

import (
    "fmt"
    "bytes"
    "strings"
    "monkey/ast"
)

const (
    IntType    = "INTEGER_TYPE"
    BoolType   = "BOOLEAN_TYPE"
    NullType   = "NULL_TYPE"
    ReturnType = "RETURN_TYPE"
    ErrorType  = "ERROR_TYPE"
    FuncType   = "FUNCTION_TYPE"
)

type ObjectType string

type Object interface {
    Type() ObjectType
    Inspect() string
    IsType(ObjectType) bool
}

type Integer struct {
    Value int64
}

// @Impl
func (this *Integer) Inspect() string {
    return fmt.Sprintf("%d", this.Value)
}

// @Impl
func (this *Integer) Type() ObjectType {
    return IntType
}

// @Impl
func (this *Integer) IsType(check ObjectType) bool {
    return check == IntType
}

type Boolean struct {
    Value bool
}

// @Impl
func (this *Boolean) Inspect() string {
    return fmt.Sprintf("%t", this.Value)
}

// @Impl
func (this *Boolean) Type() ObjectType {
    return BoolType
}

// @Impl
func (this *Boolean) IsType(check ObjectType) bool {
    return check == BoolType
}

type Null struct {}

// @Impl
func (this *Null) Inspect() string {
    return "null"
}

// @Impl
func (this *Null) Type() ObjectType {
    return NullType
}

// @Impl
func (this *Null) IsType(check ObjectType) bool {
    return check == NullType
}

type ReturnValue struct {
    Value Object
}

// @Impl
func (this *ReturnValue) Inspect() string {
    return this.Value.Inspect()
}

// @Impl
func (this *ReturnValue) Type() ObjectType {
    return ReturnType
}

// @Impl
func (this *ReturnValue) IsType(check ObjectType) bool {
    return check == ReturnType
}

type Error struct {
    Message string
}

// @Impl
func (this *Error) Inspect() string {
    return fmt.Sprintf("ERROR: %s", this.Message)
}

// @Impl
func (this *Error) Type() ObjectType {
    return ErrorType
}

// @Impl
func (this *Error) IsType(check ObjectType) bool {
    return check == ErrorType
}

type Function struct {
    Parameters []ast.Identifier
    Body *ast.StatementsBlock
    Env *Environment
}

// @Impl
func (this *Function) Inspect() string {
    var out bytes.Buffer

    var args = []string {}
    for _, arg := range this.Parameters {
        args = append(args, arg.String())
    }

    out.WriteString("fn (")
    if len(args) > 0 {
        out.WriteString(strings.Join(args, ", "))
    }
    out.WriteString(") {")
    for _, stm := range this.Body.Statements {
        out.WriteString(" " + stm.String() + "; ")
    }
    out.WriteString("}")

    return out.String()
}

// @Impl
func (this *Function) Type() ObjectType {
    return FuncType
}

// @Impl
func (this *Function) IsType(check ObjectType) bool {
    return check == FuncType
}
