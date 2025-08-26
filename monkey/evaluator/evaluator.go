// monkey/evaluator/evaluator.go

package evaluator

import (
    "monkey/object"
    "monkey/ast"
)

var (
    TRUE = &object.Boolean { Value: true }
    FALSE = &object.Boolean { Value: false }
    NULL = &object.Null {}
)

func evalStatements(statements []ast.Statement) object.Object {
    if len(statements) == 0 { return nil }
    // var evtds = []object.Object {}
    // for _, stm := range statements {
    //     var evtd = Eval(stm)
    //     if evtd != nil {
    //         evtds = append(evtds, evtd)
    //     }
    // }
    // fmt.Println(evtds)
    return Eval(statements[len(statements) - 1]) // TODO: Just takes the last to compile
}

func Eval(node ast.Node) object.Object {
    switch node := node.(type) {

// Statements
    case *ast.Program:
        return evalStatements(node.Statements)

    case *ast.ExpressionStatement:
        return Eval(node.Expression)

// Expressions
    case *ast.PrefixExpression:
        switch node.Operator {
        case "-":
            var evaluated = Eval(node.Value)
            switch x := evaluated.(type) {
            case *object.Integer:
                return &object.Integer { Value: -x.Value }
            default:
                return NULL
            }
        case "!":
            var evaluated = Eval(node.Value)
            switch x := evaluated.(type) {
            case *object.Boolean:
                if x.Value == true {
                    return FALSE
                } else {
                    return TRUE
                }
            case *object.Integer:
                if x.Value <= 0 {
                    return TRUE
                } else {
                    return FALSE
                }
            default:
                return NULL
            }
        }

    case *ast.IntegerLiteral:
        return &object.Integer { Value: node.Value }

    case *ast.Boolean:
        if node.Value {
            return TRUE
        } else {
            return FALSE
        }

    }

    return nil
}
