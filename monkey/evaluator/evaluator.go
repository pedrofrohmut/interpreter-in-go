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

func boolToObjBoolean(check bool) *object.Boolean {
    if check { return TRUE } else { return FALSE }
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
                return boolToObjBoolean(x.Value == false)
            case *object.Integer:
                return boolToObjBoolean(x.Value <= 0)
            default:
                return NULL
            }
        }

    case *ast.InfixExpression:
        var left, okLeft = Eval(node.Left).(*object.Integer)
        var right, okRight = Eval(node.Right).(*object.Integer)

        if !okLeft || !okRight {
            return NULL
        }

        switch node.Operator {
    // Math operations
        case "+":
            return &object.Integer { Value: left.Value + right.Value }
        case "-":
            return &object.Integer { Value: left.Value - right.Value }
        case "*":
            return &object.Integer { Value: left.Value * right.Value }
        case "/":
            return &object.Integer { Value: left.Value / right.Value }
    // Booleans operations
        case "==":
            return boolToObjBoolean(left.Value == right.Value)
        case "!=":
            return boolToObjBoolean(left.Value != right.Value)
        case "<":
            return boolToObjBoolean(left.Value < right.Value)
        case ">":
            return boolToObjBoolean(left.Value > right.Value)
        }


    case *ast.IntegerLiteral:
        return &object.Integer { Value: node.Value }

    case *ast.Boolean:
        return boolToObjBoolean(node.Value)

    }

    return nil
}
