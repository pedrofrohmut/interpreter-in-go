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

func isOfType(obj object.Object, objType object.ObjectType) bool {
    return obj.Type() == objType
}

func isError(obj object.Object) bool {
    return obj != nil && isOfType(obj, object.ErrorType)
}

func unwrapReturn(obj object.Object) object.Object {
    if isOfType(obj, object.ReturnType) {
        return obj.(*object.ReturnValue).Value
    }
    return obj
}

func GetMsgTypeFor(objType object.ObjectType) string {
    switch objType {
    case object.IntType:
        return "Integer"
    case object.BoolType:
        return "Boolean"
    case object.NullType:
        return "Null"
    case object.StringType:
        return "String"
    default:
        return "Not Covered"
    }
}

func getMismatchError(left object.Object, operator string, right object.Object) *object.Error {
    var leftMsgType = GetMsgTypeFor(left.Type())
    var rightMsgType = GetMsgTypeFor(right.Type())
    return &object.Error {
        Message: fmt.Sprintf("type mismatch: %s %s %s", leftMsgType, operator, rightMsgType),
    }
}

func getUnknownOperatorError(left object.Object, operator string, right object.Object) *object.Error {
    var rightMsgType = GetMsgTypeFor(right.Type())
    if left == nil {
        return &object.Error {
            Message: fmt.Sprintf("unknown operator: %s%s", operator, rightMsgType),
        }
    }
    var leftMsgType = GetMsgTypeFor(left.Type())
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

        if isError(result) { return result }

        if isOfType(result, object.ReturnType) { return result }
    }

    return result // Return the last evaluated statement if no early return types are found
}

func evalCallExpression(objFunc *object.Function, node *ast.CallExpression, env *object.Environment) object.Object {
    if len(objFunc.Parameters) != len(node.Parameters) {
        return &object.Error {
            Message: fmt.Sprintf("Expected function call to have %d parameters but found %d instead",
            len(objFunc.Parameters), len(node.Parameters)),
        }
    }

    var funcEnv = object.NewEnclosedEnvironment(objFunc.Env)

    // Adds params with values to function env
    for i := range objFunc.Parameters {
        var paramName = objFunc.Parameters[i].Value
        var paramValue = Eval(node.Parameters[i], env)
        if isError(paramValue) { return paramValue }
        funcEnv.Set(paramName, paramValue)
    }

    var result = Eval(objFunc.Body, funcEnv)

    return unwrapReturn(result)
}

func findIdentifier(name string, env *object.Environment) object.Object {
    var value, ok = env.Get(name)
    if ok { return value }

    var builtin, okBuiltin = GetBuiltin(name)
    if okBuiltin { return builtin }

    return getIdentifierNotFoundError(name)
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

    case *ast.StatementsBlock:
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

        if evaluatedLeft.Type() != evaluatedRight.Type() {
            return getMismatchError(evaluatedLeft, node.Operator, evaluatedRight)
        }

        switch evaluatedLeft.Type() {
        case object.StringType:
            if node.Operator != "+" {
                return getUnknownOperatorError(evaluatedLeft, node.Operator, evaluatedRight)
            }
            var left = evaluatedLeft.(*object.String)
            var right = evaluatedRight.(*object.String)
            return &object.String { Value: left.Value + right.Value }
        case object.IntType:
            var left = evaluatedLeft.(*object.Integer)
            var right = evaluatedRight.(*object.Integer)
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
        case object.BoolType:
            return getUnknownOperatorError(evaluatedLeft, node.Operator, evaluatedRight)
        default:
            return getNotCoveredEvaluationError(node)
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

    case *ast.FunctionLiteral:
        return &object.Function { Parameters: node.Parameters, Body: node.Body, Env: env }

    case *ast.CallExpression:
        switch exp := node.Expression.(type) {
        case *ast.Identifier: // Exp: foo(x, y, z)
            var obj = findIdentifier(exp.Value, env)
            if isError(obj) { return obj }

            switch objFunc := obj.(type) {
            case *object.Function:
                return evalCallExpression(objFunc, node, env)
            case *object.Builtin:
                var args = []object.Object {}
                for _, param := range node.Parameters {
                    var arg = Eval(param, env)
                    if isError(arg) { return arg }
                    args = append(args, arg)
                }
                var fn = objFunc.Function
                return fn(args...)
            default:
                return &object.Error {
                    Message: fmt.Sprintf("Identifier is not connected to an covered function type. Found %T instead", obj),
                }
            }

        case *ast.FunctionLiteral: // Exp: fn (x, y) { x + y; }(5, 6)
            var obj = Eval(exp, env)
            if isError(obj) { return obj }

            var objFunc = obj.(*object.Function)

            return evalCallExpression(objFunc, node, env)

        default:
            return &object.Error {
                Message:fmt.Sprintf("Not covered CallExpression.Expression type: %T", node.Expression),
            }
        }

    case *ast.Identifier:
        return findIdentifier(node.Value, env)

    case *ast.IntegerLiteral:
        return &object.Integer { Value: node.Value }

    case *ast.Boolean:
        return objFromBool(node.Value)

    case *ast.StringLiteral:
        return &object.String { Value: node.Value }

    } // END: Switch nody.(type)

    return getNotCoveredEvaluationError(node)
}
