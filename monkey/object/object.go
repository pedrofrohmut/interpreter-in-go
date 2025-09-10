// monkey/object/object.go

package object

import (
    "bytes"
    "fmt"
    "monkey/ast"
    "strings"
    "hash/fnv"
)

const (
    IntType     = "INTEGER_TYPE"
    BoolType    = "BOOLEAN_TYPE"
    NullType    = "NULL_TYPE"
    ReturnType  = "RETURN_TYPE"
    ErrorType   = "ERROR_TYPE"
    FuncType    = "FUNCTION_TYPE"
    StringType  = "STRING_TYPE"
    BuiltinType = "BUILTIN_TYPE"
    ArrayType   = "ARRAY_TYPE"
    CharType    = "CHAR_TYPE"
    HashType    = "HASH_TYPE"
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

type Char struct {
    Value byte
}

// @Impl
func (this *Char) Inspect() string {
    return string(this.Value)
}

// @Impl
func (this *Char) Type() ObjectType {
    return CharType
}

type String struct {
    Value string
}

// @Impl
func (this *String) Inspect() string {
    return this.Value
}

// @Impl
func (this *String) Type() ObjectType {
    return StringType
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

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
    Function BuiltinFunction
}

// @Impl
func (this *Builtin) Inspect() string {
    return "builtin function"
}

// @Impl
func (this *Builtin) Type() ObjectType {
    return BuiltinType
}

type Array struct {
    Elements []Object
}

// @Impl
func (this *Array) Inspect() string {
    var out bytes.Buffer

    var elements = []string {}
    for _, element := range this.Elements {
        elements = append(elements, element.Inspect())
    }

    out.WriteString("[")
    out.WriteString(strings.Join(elements, ", "))
    out.WriteString("]")

    return out.String()
}

// @Impl
func (this *Array) Type() ObjectType {
    return ArrayType
}

type Hashable interface {
    HashKey() HashKey
}

type HashKey struct {
    Type ObjectType
    Value uint64
}

// @Impl
func (this *Boolean) HashKey() HashKey {
    var value uint64; if this.Value { value = 1 } else { value = 0 } // This crap because there is no ternary operator
    return HashKey { Type: this.Type(), Value: value }
}

// @Impl
func (this *Integer) HashKey() HashKey {
    return HashKey { Type: this.Type(), Value: uint64(this.Value) }
}

// @Impl
func (this *String) HashKey() HashKey {
    var hash = fnv.New64a()
    hash.Write([]byte(this.Value))
    var value = hash.Sum64()
    return HashKey { Type: this.Type(), Value: value }
}

// Had to create this struct because i could not keep the OriginalKey in the HashKey struct
// it breaks the comparison capabilities in golang
// Key used as structs must have only comparable types. More info: https://go.dev/blog/maps
type HashPair struct {
    OriginalKey Object
    Value Object
}

type Hash struct {
    Pairs map[HashKey]HashPair
}

// @Impl
func (this *Hash) Inspect() string {
    var out bytes.Buffer

    if len(this.Pairs) == 0 { return "{}" }

    var pairs = []string {}
    for _, value := range this.Pairs {
        var pair = value.OriginalKey.Inspect() + ": " + value.Value.Inspect()
        pairs = append(pairs, pair)
    }

    out.WriteString("{ ")
    out.WriteString(strings.Join(pairs, ", "))
    out.WriteString(" }")

    return out.String()
}

// @Impl
func (this *Hash) Type() ObjectType {
    return HashType
}
