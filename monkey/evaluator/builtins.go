// monkey/evaluator/builtins.go

package evaluator

import (
    "monkey/object"
    "fmt"
    "strings"
)

func getNumArgsError(expected int, got int) *object.Error {
    return &object.Error {
        Message: fmt.Sprintf("wrong number of arguments. expected=%d but got=%d", expected, got),
    }
}

func getTypeNotSupportedError(funcName string, obj object.Object) *object.Error {
    var unsupportedType = GetMsgTypeFor(obj.Type())
    return &object.Error {
        Message: fmt.Sprintf("argument to %s not supported, got %s", funcName, unsupportedType),
    }
}

// Returns the length of an array or string
var Len = func (args ...object.Object) object.Object {
    if len(args) != 1 {
        return getNumArgsError(1, len(args))
    }

    var first = args[0]

    switch obj := first.(type) {
    case *object.String:
        return &object.Integer { Value: int64(len(obj.Value)) }
    case *object.Array:
        return &object.Integer { Value: int64(len(obj.Elements)) }
    default:
        return getTypeNotSupportedError("len", obj)
    }
}

// Returns the first element of an array or string
var First = func (args ...object.Object) object.Object {
    if len(args) != 1 {
        return getNumArgsError(1, len(args))
    }

    switch obj := args[0].(type) {
    case *object.String:
        if len(obj.Value) == 0 {
            return ObjNull
        }
        var ch = []byte(obj.Value)[0]
        return &object.Char { Value: ch }
    case *object.Array:
        if len(obj.Elements) == 0 {
            return ObjNull
        }
        return obj.Elements[0]
    default:
        return getTypeNotSupportedError("first", obj)
    }
}

// Returns the last element of an array or a string
var Last = func (args ...object.Object) object.Object {
    if len(args) != 1 {
        return getNumArgsError(1, len(args))
    }

    switch obj := args[0].(type) {
    case *object.String:
        if len(obj.Value) == 0 {
            return ObjNull
        }
        var ch = []byte(obj.Value)[len(obj.Value) - 1]
        return &object.Char { Value: ch }
    case *object.Array:
        if len(obj.Elements) == 0 {
            return ObjNull
        }
        return obj.Elements[len(obj.Elements) - 1]
    default:
        return getTypeNotSupportedError("last", obj)
    }
}

// Returns a copy from the source array without the first element
var Rest = func (args ...object.Object) object.Object {
    if len(args) != 1 {
        return getNumArgsError(1, len(args))
    }

    switch obj := args[0].(type) {
    case *object.String:
        if len(obj.Value) < 2 { // Return empty string when too small
            return &object.String { Value: "" }
        }
        var copy = strings.Clone(obj.Value)
        var rest = copy[1:]
        return &object.String { Value: rest }
    case *object.Array:
        if len(obj.Elements) < 2 { // Return empty array when is too small
            var arr = &object.Array {}
            arr.Elements = []object.Object {}
            return arr
        }
        var rest = make([]object.Object, len(obj.Elements) - 1)
        copy(rest, obj.Elements[1:])
        return &object.Array { Elements: rest }
    default:
        return getTypeNotSupportedError("rest", obj)
    }
}

var Push = func (args ...object.Object) object.Object {
    if len(args) != 2 {
        return getNumArgsError(2, len(args))
    }

    var arr, okArr = args[0].(*object.Array)
    if !okArr {
        return &object.Error { Message: "First argument of push function must be an array" }
    }

    arr.Elements = append(arr.Elements, args[1]);
    return arr
}

var Puts = func (args ...object.Object) object.Object {
    if len(args) != 1 {
        return getNumArgsError(2, len(args))
    }

    fmt.Println(args[0].Inspect())

    return ObjNull
}

func GetBuiltin(name string) (*object.Builtin, bool) {
    switch name {
    case "len":
        return &object.Builtin { Function: Len }, true
    case "first":
        return &object.Builtin { Function: First }, true
    case "last":
        return &object.Builtin { Function: Last }, true
    case "rest":
        return &object.Builtin { Function: Rest }, true
    case "push":
        return &object.Builtin { Function: Push }, true
    case "puts":
        return &object.Builtin { Function: Puts }, true
    }

    return nil, false
}
