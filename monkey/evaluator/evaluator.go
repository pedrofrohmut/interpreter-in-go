// monkey/evaluator/evaluator.go

package evaluator

import (
    "monkey/object"
    "monkey/ast"
    _"fmt"
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
    case *ast.IntegerLiteral:
        return &object.Integer { Value: node.Value }

    case *ast.Boolean:
        return &object.Boolean { Value: node.Value }

    }

    return nil
}
