// monkey/evaluator/evaluator_test.go

package evaluator

import (
    _"fmt"
    "testing"
    "monkey/lexer"
    "monkey/parser"
    "monkey/utils"
    "monkey/object"
)

func TestEvalIntegerExpression(t *testing.T) {
    var tests = []struct {
        input string; expected int64
    } {
        { "5", 5 },
        { "10", 10 },
        { "15", 15 },
    }

    var input = utils.TryGetInput(t, tests)
    var lexer = lexer.NewLexer(input)
    var parser = parser.NewParser(lexer)
    var program = parser.ParseProgram()

    for i, stm := range program.Statements {
        var evaluated = Eval(stm)
        var res, ok = evaluated.(*object.Integer)
        if !ok {
            t.Errorf("Evaluated statement was not evaluated to an object.Integer. Got %T instead", evaluated)
            continue
        }
        if res.Value != tests[i].expected {
            t.Errorf("Expected result object value to be %d but got %d instead", tests[i].expected, res.Value)
        }
    }
}
