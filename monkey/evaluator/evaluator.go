// monkey/evaluator/evaluator.go

package evaluator

import (
    "fmt"
    "monkey/object"
    "monkey/ast"
)

var (
    ObjTrue = &object.Boolean { Value: true }
    ObjFalse = &object.Boolean { Value: false }
    ObjNull = &object.Null {}
)

func getMsgTypeFor(objType object.ObjectType) string {
    switch objType {
    case object.IntType:
        return "Integer"
    case object.BoolType:
        return "Boolean"
    case object.NullType:
        return "Null"
    default:
        return "TODO: Not Covered"
    }
}

func getMismatchError(left object.Object, operator string, right object.Object) *object.Error {
    var leftMsgType = getMsgTypeFor(left.Type())
    var rightMsgType = getMsgTypeFor(right.Type())
    return &object.Error {
        Message: fmt.Sprintf("type mismatch: %s %s %s", leftMsgType, operator, rightMsgType),
    }
}

func getUnknownOperatorError(left object.Object, operator string, right object.Object) *object.Error {
    var rightMsgType = getMsgTypeFor(right.Type())
    if left == nil {
        return &object.Error {
            Message: fmt.Sprintf("unknown operator: %s%s", operator, rightMsgType),
        }
    }
    var leftMsgType = getMsgTypeFor(left.Type())
    return &object.Error {
        Message: fmt.Sprintf("unknown operator: %s %s %s", leftMsgType, operator, rightMsgType),
    }
}

func evalStatements(statements []ast.Statement) object.Object {
    if len(statements) == 0 { return nil }

    var result object.Object

    for _, stm := range statements {
        result = Eval(stm)

        if result.Type() == object.ErrorType { // Return early when error type is found
            return result
        }

        if result.Type() == object.ReturnType { // Return early when return type is found
            var x = result.(*object.ReturnValue)

            if x.Value.Type() == object.ErrorType {
                return x.Value
            }

            return result
        }
    }

    return result // Return the last evaluated statement if no early return types are found
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
            case *object.Boolean:
                return getUnknownOperatorError(nil, node.Operator, evaluated)
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
        var evaluatedLeft, evaluatedRight = Eval(node.Left), Eval(node.Right)

        var left, okLeft = evaluatedLeft.(*object.Integer)
        var right, okRight = evaluatedRight.(*object.Integer)

        if !okLeft && !okRight {
            return getUnknownOperatorError(evaluatedLeft, node.Operator, evaluatedRight)
        }

        if !okLeft || !okRight {
            return getMismatchError(evaluatedLeft, node.Operator, evaluatedRight)
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

// TODO: Make if eval all needed statements (can use evalStatements)
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
