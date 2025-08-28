// monkey/evaluator/evaluator_test.go

package evaluator

import (
    "testing"
    "monkey/lexer"
    "monkey/parser"
    "monkey/object"
    "monkey/test_utils"
)

func TestEvalIntegerExpression(t *testing.T) {
    var tests = []struct {
        input string; expected int64
    } {
        { "5", 5 },
        { "10", 10 },
        { "15", 15 },

        // Prefix Expressions
        { "-5", -5 },
        { "-10", -10 },
        { "-15", -15 },

        // Infix Expressions
        { "5 + 5", 10 },
        { "5 + 5 + 5", 15 },
        { "5 + 5 + 5 + 5 - 10", 10 },
        { "2 * 2 * 2 * 2 * 2", 32 },
        { "-50 + 100 + -50", 0 },
        { "5 * 2 + 10", 20 },
        { "5 + 2 * 10", 25 },
        { "20 + 2 * -10", 0 },
        { "50 / 2 * 2 + 10", 60 },
        { "2 * (5 + 10)", 30 },
        { "3 * 3 * 3 + 10", 37 },
        { "3 * (3 * 3) + 10", 37 },
        { "(5 + 10 * 2 + 15 / 3) * 2 + -10", 50 },
    }

    var input = test_utils.TryGetInput(t, tests)
    var lexer = lexer.NewLexer(input)
    var parser = parser.NewParser(lexer)
    var program = parser.ParseProgram()

    test_utils.CheckForParserErrors(t, parser)
    test_utils.CheckProgram(t, program, len(tests))

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

func TestEvalBooleanExpression(t *testing.T) {
    var tests = []struct {
        input string; expected bool
    } {
        { "true",  true  },
        { "false", false },

        // Infix Expressions
        { "1 < 2",  true  },
        { "1 > 2",  false },
        { "1 < 1",  false },
        { "1 > 1",  false },
        { "1 == 1", true  },
        { "1 != 1", false },
        { "1 == 2", false },
        { "1 != 2", true  },
    }

    var input = test_utils.TryGetInput(t, tests)
    var lexer = lexer.NewLexer(input)
    var parser = parser.NewParser(lexer)
    var program = parser.ParseProgram()

    test_utils.CheckForParserErrors(t, parser)
    test_utils.CheckProgram(t, program, len(tests))

    for i, stm := range program.Statements {
        var evaluated = Eval(stm)
        var res, ok = evaluated.(*object.Boolean)
        if !ok {
            t.Errorf("Evaluated statement was not evaluated to an object.Boolean. Got %T instead", evaluated)
            continue
        }
        if res.Value != tests[i].expected {
            t.Errorf("Expected result object value to be %t but got %t instead", tests[i].expected, res.Value)
        }
    }
}

func TestEvalBangOperator(t *testing.T) {
    var tests = []struct {
        input string; expected bool
    } {
        { "!true", false },
        { "!false", true },
        { "!!true", true },
        { "!!false", false },
        { "!5", false },
        { "!!5", true },
    }

    var input = test_utils.TryGetInput(t, tests)
    var lexer = lexer.NewLexer(input)
    var parser = parser.NewParser(lexer)
    var program = parser.ParseProgram()

    test_utils.CheckForParserErrors(t, parser)
    test_utils.CheckProgram(t, program, len(tests))

    for i, stm := range program.Statements {
        var evaluated = Eval(stm)
        var res, ok = evaluated.(*object.Boolean)
        if !ok {
            t.Errorf("Evaluated statement was not evaluated to an object.Boolean. Got %T instead", evaluated)
            continue
        }
        if res.Value != tests[i].expected {
            t.Errorf("Expected result object value to be %t but got %t instead", tests[i].expected, res.Value)
        }
    }
}

func TestIfElseExpressions(t *testing.T) {
    var tests = []struct {
        input string; expected any
    } {
        { "if (true) { 10 }", 10 },
        { "if (false) { 10 }", nil },
        { "if (1) { 10 }", 10 },
        { "if (1 < 2) { 10 }", 10 },
        { "if (1 > 2) { 10 }", nil },
        { "if (1 > 2) { 10 } else { 20 }", 20 },
        { "if (1 < 2) { 10 } else { 20 }", 10 },

        { "if (true) { 123 } else { 666 }", 123 },
        { "if (false) { 123 } else { 666 }", 666 },
    }

    var input = test_utils.TryGetInput(t, tests)
    var lexer = lexer.NewLexer(input)
    var parser = parser.NewParser(lexer)
    var program = parser.ParseProgram()

    test_utils.CheckForParserErrors(t, parser)
    test_utils.CheckProgram(t, program, len(tests))

    for i, stm := range program.Statements {
        var evaluated = Eval(stm)

        switch x := tests[i].expected.(type) { // Switch on expected type
        case int:
            var res, ok = evaluated.(*object.Integer)
            if !ok {
                t.Errorf("Expected evaluated object to be type of object.Integer. Got %T instead", evaluated)
                continue
            }
            if res.Value != int64(x) {
                t.Errorf("Expected result object.Integer value to be %d but got %d instead", x, res.Value)
            }
        case nil:
           var _, ok = evaluated.(*object.Null)
           if !ok {
               t.Errorf("Expected evaluated object to be type of object.Null. Got %T instead", evaluated)
               continue
           }
        }
    }
}

func TestReturnStatements(t *testing.T) {
    var tests = []struct {
        input string; expected int64
    } {
        { "return 5", 5 },
        { "return 10", 10 },
        { "return 15", 15 },
        { "return 2 * 10", 20 },
        { "return 2 * 5; 9", 10 },
        { "9; return 3 * 7; 9", 21 },
        {
            `
                if (10 > 1) {
                    if (10 > 1) {
                        return 10;
                    }
                    return 1;
                }
            `,
            10,
        },
    }

    for _, test := range tests {
        var lexer = lexer.NewLexer(test.input)
        var parser = parser.NewParser(lexer)
        var program = parser.ParseProgram()

        test_utils.CheckForParserErrors(t, parser)

        var evaluated = Eval(program)
        var returnObj, ok = evaluated.(*object.ReturnValue)
        if !ok {
            t.Errorf("Expected evaluated object to be type of object.ReturnValue. Got %T instead", evaluated)
            continue
        }
        var intObj, ok2 = returnObj.Value.(*object.Integer)
        if !ok2 {
            t.Errorf("Expected return object value to be type of object.Integer. Got %T instead", returnObj.Value)
            continue
        }
        if intObj.Value != test.expected {
            t.Errorf("Expected return value object to be object.Integer with value %d but got %d instead", test.expected, intObj.Value)
        }
    }
}
