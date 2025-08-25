// monkey/object/object.go

package object

import (
    "fmt"
)

const (
    INT_OBJ = "INTEGER"
    BOOL_OBJ = "BOOLEAN"
    NULL_OBJ = "NULL"
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
    return INT_OBJ
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
    return BOOL_OBJ
}

type Null struct {}

// @Impl
func (this *Null) Inspect() string {
    return "null"
}

// @Impl
func (this *Null) Type() ObjectType {
    return NULL_OBJ
}
