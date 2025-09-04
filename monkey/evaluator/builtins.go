// monkey/evaluator/builtins.go

package evaluator

import (
    "monkey/object"
    "fmt"
)

var Len = func (args ...object.Object) object.Object {
    var first = args[0]

    if len(args) > 1 {
        return &object.Error {
            Message: fmt.Sprintf("wrong number of arguments. expected=%d but got=%d", 1, len(args)),
        }
    }

    switch obj := first.(type) {
    case *object.String:
        return &object.Integer { Value: int64(len(obj.Value)) }
    default:
        return &object.Error {
            Message: fmt.Sprintf("argument to len not supported, got %s", GetMsgTypeFor(obj.Type())),
        }
    }
}

func GetBuiltin(name string) (*object.Builtin, bool) {
    switch name {
    case "len":
        return &object.Builtin { Function: Len }, true
    }

    return nil, false
}
