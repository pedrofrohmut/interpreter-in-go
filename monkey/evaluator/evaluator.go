// monkey/evaluator/evaluator.go

package evaluator

import (
    "fmt"
    "monkey/object"
    "monkey/ast"
)

var (
    ObjTrue  = &object.Boolean { Value: true }
    ObjFalse = &object.Boolean { Value: false }
    ObjNull  = &object.Null {}
)

func isError(obj object.Object) bool {
    return obj != nil && obj.Type() == object.ErrorType
}

func getMsgTypeFor(objType object.ObjectType) string {
    switch objType {
    case object.IntType:
        return "Integer"
    case object.BoolType:
        return "Boolean"
    case object.NullType:
        return "Null"
    default:
        return "Not Covered"
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

func getNotCoveredEvaluationError(node ast.Node) *object.Error {
    return &object.Error {
        Message: fmt.Sprintf("Node %T not covered in evaluation", node),
    }
}

func getIdentifierNotFoundError(name string) *object.Error {
    return &object.Error {
        Message: fmt.Sprintf("identifier not found: %s", name),
    }
}

func evalStatements(statements []ast.Statement, env *object.Environment) object.Object {
    var result object.Object = nil

    for _, stm := range statements {
        result = Eval(stm, env)

        // Return early when error type is found
        if isError(result) { return result }

        // Return early when return type is found
        if result.IsType(object.ReturnType) { return result }
    }

    // Return the last evaluated statement if no early return types are found
    return result
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

func Eval(node ast.Node, env *object.Environment) object.Object {
    switch node := node.(type) {

// Statements
    case *ast.Program:
        return evalStatements(node.Statements, env)

    case *ast.ReturnStatement:
        var value = Eval(node.Expression, env)
        if isError(value) { return value }

        switch x := value.(type) {
        case *object.Null:
            return &object.ReturnValue { Value: ObjNull }
        default:
            return &object.ReturnValue { Value: x }
        }

    case *ast.LetStatement:
        var expValue = Eval(node.Expression, env)
        if isError(expValue) { return expValue }

        env.Set(node.Identifier, expValue)

        return ObjNull

    case *ast.ExpressionStatement:
        return Eval(node.Expression, env)

// Expressions
    case *ast.PrefixExpression:
        switch node.Operator {
        case "-":
            var evaluated = Eval(node.Value, env)
            if isError(evaluated) { return evaluated }

            switch x := evaluated.(type) {
            case *object.Integer:
                return &object.Integer { Value: -x.Value }
            default:
                return getUnknownOperatorError(nil, node.Operator, evaluated)
            }
        case "!":
            var evaluated = Eval(node.Value, env)
            if isError(evaluated) { return evaluated }

            switch x := evaluated.(type) {
            case *object.Boolean:
                return objFromBool(!isTruthy(x.Value))
            case *object.Integer:
                return objFromBool(!isTruthy(x.Value))
            default:
                return getUnknownOperatorError(nil, node.Operator, evaluated)
            }
        }

    case *ast.InfixExpression:
        var evaluatedLeft = Eval(node.Left, env)
        var evaluatedRight =  Eval(node.Right, env)

        if isError(evaluatedLeft) { return evaluatedLeft }
        if isError(evaluatedRight) { return evaluatedRight }

        if evaluatedLeft.Type() != evaluatedRight.Type() { // Types are different
            return getMismatchError(evaluatedLeft, node.Operator, evaluatedRight)
        }

        var left, okLeft = evaluatedLeft.(*object.Integer)
        var right, okRight = evaluatedRight.(*object.Integer)

        if !okLeft && !okRight {
            return getUnknownOperatorError(evaluatedLeft, node.Operator, evaluatedRight)
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

// TODO: Make IfExpression good and not this mess
// TODO: Make if eval all needed statements (can use evalStatements)
    case *ast.IfExpression:
        var conditionResult = Eval(node.Condition, env)
        if isError(conditionResult) { return conditionResult }

        switch x := conditionResult.(type) {
        case *object.Boolean:
            if isTruthy(x.Value) {
                return Eval(node.ConsequenceBlock.Statements[0], env)
            } else if node.AlternativeBlock != nil {
                return Eval(node.AlternativeBlock.Statements[0], env)
            } else {
                // TODO: error for missing alternative
                return ObjNull
            }
        case *object.Integer:
            if isTruthy(x.Value) {
                return Eval(node.ConsequenceBlock.Statements[0], env)
            } else if node.AlternativeBlock != nil {
                return Eval(node.AlternativeBlock.Statements[0], env)
            } else {
                // TODO: error for missing alternative
                return ObjNull
            }
        }

    case *ast.Identifier:
        var val, ok = env.Get(node.Value)
        if !ok {
            return getIdentifierNotFoundError(node.Value)
        }
        return val

    case *ast.IntegerLiteral:
        return &object.Integer { Value: node.Value }

    case *ast.Boolean:
        return objFromBool(node.Value)

    } // END: Switch nody.(type)

    return getNotCoveredEvaluationError(node)
}
