// monkey/object/object.go

package object

import (
    "fmt"
)

const (
    IntType    = "INTEGER_TYPE"
    BoolType   = "BOOLEAN_TYPE"
    NullType   = "NULL_TYPE"
    ReturnType = "RETURN_TYPE"
    ErrorType  = "ERROR_TYPE"
)

type ObjectType string

type Object interface {
    Type() ObjectType
    Inspect() string
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

type Null struct {}

// @Impl
func (this *Null) Inspect() string {
    return "null"
}

// @Impl
func (this *Null) Type() ObjectType {
    return NullType
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
