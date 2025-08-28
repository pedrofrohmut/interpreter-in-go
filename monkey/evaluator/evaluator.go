// monkey/evaluator/evaluator.go

package evaluator

import (
    "monkey/object"
    "monkey/ast"
)

var (
    ObjTrue = &object.Boolean { Value: true }
    ObjFalse = &object.Boolean { Value: false }
    ObjNull = &object.Null {}
)

func evalStatements(statements []ast.Statement) object.Object {
    if len(statements) == 0 { return nil }

    var result object.Object

    for _, stm := range statements {
        result = Eval(stm)
        if result.Type() == object.ReturnType { // Return early when return type is found
            return result
        }
    }

    return result // Return the last eval if no return is defined
}

func objFromBool(check bool) *object.Boolean {
    if check { return ObjTrue } else { return ObjFalse }
}

func isTruthy(check any) bool {
    switch x := check.(type) {
    case int:
        return x > 0
    case int64:
        return x > 0
    case bool:
        return x
    default:
        return false
    }
}

func Eval(node ast.Node) object.Object {
    switch node := node.(type) {

// Statements
    case *ast.Program:
        return evalStatements(node.Statements)

    case *ast.ReturnStatement:
        var value = Eval(node.Expression)
        switch x := value.(type) {
        case *object.Null:
            return &object.ReturnValue { Value: ObjNull }
        default:
            return &object.ReturnValue { Value: x }
        }

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
                return ObjNull
            }
        case "!":
            var evaluated = Eval(node.Value)
            switch x := evaluated.(type) {
            case *object.Boolean:
                return objFromBool(!isTruthy(x.Value))
            case *object.Integer:
                return objFromBool(!isTruthy(x.Value))
            default:
                return ObjNull
            }
        }

    case *ast.InfixExpression:
        var left, okLeft = Eval(node.Left).(*object.Integer)
        var right, okRight = Eval(node.Right).(*object.Integer)
        if !okLeft || !okRight {
            return ObjNull
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
            return objFromBool(left.Value == right.Value)
        case "!=":
            return objFromBool(left.Value != right.Value)
        case "<":
            return objFromBool(left.Value < right.Value)
        case ">":
            return objFromBool(left.Value > right.Value)
        }

// TODO: Make if eval all needed statements
    case *ast.IfExpression:
        var conditionResult = Eval(node.Condition)
        switch x := conditionResult.(type) {
        case *object.Boolean:
            if isTruthy(x.Value) {
                return Eval(node.ConsequenceBlock.Statements[0])
            } else if node.AlternativeBlock != nil {
                return Eval(node.AlternativeBlock.Statements[0])
            } else {
                return ObjNull
            }
        case *object.Integer:
            if isTruthy(x.Value) {
                return Eval(node.ConsequenceBlock.Statements[0])
            } else if node.AlternativeBlock != nil {
                return Eval(node.AlternativeBlock.Statements[0])
            } else {
                return ObjNull
            }
        }

    case *ast.IntegerLiteral:
        return &object.Integer { Value: node.Value }

    case *ast.Boolean:
        return objFromBool(node.Value)

    }

    return nil
}
