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
        { "5",  5  },
        { "10", 10 },
        { "15", 15 },

        // Prefix Expressions
        { "-5",  -5  },
        { "-10", -10 },
        { "-15", -15 },

        // Infix Expressions
        { "5 + 5",                           10 },
        { "5 + 5 + 5",                       15 },
        { "5 + 5 + 5 + 5 - 10",              10 },
        { "2 * 2 * 2 * 2 * 2",               32 },
        { "-50 + 100 + -50",                 0  },
        { "5 * 2 + 10",                      20 },
        { "5 + 2 * 10",                      25 },
        { "20 + 2 * -10",                    0  },
        { "50 / 2 * 2 + 10",                 60 },
        { "2 * (5 + 10)",                    30 },
        { "3 * 3 * 3 + 10",                  37 },
        { "3 * (3 * 3) + 10",                37 },
        { "(5 + 10 * 2 + 15 / 3) * 2 + -10", 50 },
    }

    var input = test_utils.TryGetInput(t, tests)
    var lexer = lexer.NewLexer(input)
    var parser = parser.NewParser(lexer)
    var program = parser.ParseProgram()

    test_utils.CheckForParserErrors(t, parser)
    test_utils.CheckProgram(t, program, len(tests))

    for i, stm := range program.Statements {
        var env = object.NewEnvironment()
        var evaluated = Eval(stm, env)
        if (test_utils.CheckForEvalError(t, evaluated)) { continue }

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
        var env = object.NewEnvironment()
        var evaluated = Eval(stm, env)
        if (test_utils.CheckForEvalError(t, evaluated)) { continue }

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
        { "!true",   false },
        { "!false",  true  },
        { "!!true",  true  },
        { "!!false", false },
        { "!5",      false },
        { "!!5",     true  },
    }

    var input = test_utils.TryGetInput(t, tests)
    var lexer = lexer.NewLexer(input)
    var parser = parser.NewParser(lexer)
    var program = parser.ParseProgram()

    test_utils.CheckForParserErrors(t, parser)
    test_utils.CheckProgram(t, program, len(tests))

    for i, stm := range program.Statements {
        var env = object.NewEnvironment()
        var evaluated = Eval(stm, env)
        if (test_utils.CheckForEvalError(t, evaluated)) { continue }

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
        { "if (true) { 10 }",                10  },
        { "if (false) { 10 }",               nil },
        { "if (1) { 10 }",                   10  },
        { "if (1 < 2) { 10 }",               10  },
        { "if (1 > 2) { 10 }",               nil },
        { "if (1 > 2) { 10 } else { 20 }",   20  },
        { "if (1 < 2) { 10 } else { 20 }",   10  },

        { "if (true) { 123 } else { 666 }",  123 },
        { "if (false) { 123 } else { 666 }", 666 },

        // TODO: Add more tests where the consequence and alternative blocks have
        // more stuff like assign variables and return statements
    }

    var input = test_utils.TryGetInput(t, tests)
    var lexer = lexer.NewLexer(input)
    var parser = parser.NewParser(lexer)
    var program = parser.ParseProgram()

    test_utils.CheckForParserErrors(t, parser)
    test_utils.CheckProgram(t, program, len(tests))

    for i, stm := range program.Statements {
        var env = object.NewEnvironment()
        var evaluated = Eval(stm, env)
        if (test_utils.CheckForEvalError(t, evaluated)) { continue }

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
        { "return 5;",           5  },
        { "return 10;",          10 },
        { "return 15;",          15 },
        { "return 2 * 10;",      20 },
        { "return 2 * 5; 9;",    10 },
        { "9; return 3 * 7; 9;", 21 },
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

        var env = object.NewEnvironment()
        var evaluated = Eval(program, env)
        if (test_utils.CheckForEvalError(t, evaluated)) { continue }

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

func TestErrorHandling(t *testing.T) {
    var tests = []struct {
        input string; expected string
    } {
        { "5 + true;",                     "type mismatch: Integer + Boolean"    },
        { "5 + true; 9;",                  "type mismatch: Integer + Boolean"    },
        { "-true;",                        "unknown operator: -Boolean"          },
        { "true + false;",                 "unknown operator: Boolean + Boolean" },
        { "5; true + false; 5;",           "unknown operator: Boolean + Boolean" },
        { "if (10 > 1) { true + false; }", "unknown operator: Boolean + Boolean" },
        {
            `
                if (10 > 1) {
                    if (10 > 1) {
                        return true + false;
                    }
                    return 1;
                }
            `,
            "unknown operator: Boolean + Boolean",
        },

        // Let Statements
        { "foobar;", "identifier not found: foobar" },
    }

    for _, test := range tests {
        var lexer = lexer.NewLexer(test.input)
        var parser = parser.NewParser(lexer)
        var program = parser.ParseProgram()

        test_utils.CheckForParserErrors(t, parser)

        var env = object.NewEnvironment()
        var evaluated = Eval(program, env)
        var errObj, ok = evaluated.(*object.Error)
        if !ok {
            t.Errorf("Expected evaluated object to be of object.Error. Got %T instead", evaluated)
            continue
        }
        if errObj.Message != test.expected {
            t.Errorf("Expected error object message to be '%s' but got '%s' instead", test.expected, errObj.Message)
        }
    }
}

func TestLetStatements(t *testing.T) {
    var tests = []struct {
        input string; expected int64
    } {
        { "let a = 5; a;", 5 },
        { "let a = 5; let b = a; b;", 5 },
        { "let a = 5 * 5; a;", 25 },
        { "let a = 5; let b = 10; let c = a + b; c;", 15 },
    }

    for _, test := range tests {
        var lexer = lexer.NewLexer(test.input)
        var parser = parser.NewParser(lexer)
        var program = parser.ParseProgram()

        test_utils.CheckForParserErrors(t, parser)

        var env = object.NewEnvironment()
        var evaluated = Eval(program, env)
        if (test_utils.CheckForEvalError(t, evaluated)) { continue }

        var val, ok = evaluated.(*object.Integer)
        if !ok {
            t.Errorf("Expected evaluated object to be type of object.Integer. Got %T instead", evaluated)
            continue
        }
        if val.Value != test.expected {
            t.Errorf("Expected evaluated object value to be %d but got %d instead", test.expected, val.Value)
        }
    }
}

func TestFunctionObject(t *testing.T) {
    var input = "fn (x) { x + 2; };"

    var lexer = lexer.NewLexer(input)
    var parser = parser.NewParser(lexer)
    var program = parser.ParseProgram()

    test_utils.CheckForParserErrors(t, parser)

    var env = object.NewEnvironment()
    var evaluated = Eval(program, env)
    if (test_utils.CheckForEvalError(t, evaluated)) { return }

    var val, ok = evaluated.(*object.Function)
    if !ok {
        t.Errorf("Expected evaluated object to be type of object.Function. Got %T instead", evaluated)
        return
    }
    // Arguments
    if val.Parameters[0].Value != "x" {
        t.Errorf("Expected function first argument to be %s but got %s instead", "x", val.Parameters[0].Value)
    }
    // Body
    var expectedBody = "(x + 2)"
    if val.Body.String() != expectedBody {
        t.Errorf("Expected function body to be '%s' but got '%s' instead", expectedBody, val.Body.String())
    }
}

func TestFunctionApplication(t *testing.T) {
    var tests = []struct {
        input string; expected int64
    } {
        { "let identity = fn(x) { x; }; identity(5);",                      5  },
        { "let identity = fn(x) { return x; }; identity(5);",               5  },
        { "let double = fn(x) { x * 2; }; double(5);",                      10 },
        { "let add = fn(x, y) { x + y; }; add(10, 5);",                     15 },
        { "let add = fn(x, y) { return x + y; }; add(10, 5);",              15 },
        { "let add = fn(x, y) { return x + y; }; add(10 + 5, add(5, 10));", 30 },
        { "fn (x) { x; }(5)",                                               5  },

        // Custom
        { "let a = 5; let add = fn(x) { a + x; }; add(10);",                15 },
    }

    for _, test := range tests {
        var lexer = lexer.NewLexer(test.input)
        var parser = parser.NewParser(lexer)
        var program = parser.ParseProgram()

        test_utils.CheckForParserErrors(t, parser)

        var env = object.NewEnvironment()
        var evaluated = Eval(program, env)
        if test_utils.CheckForEvalError(t, evaluated) { continue }

        var val, ok = evaluated.(*object.Integer)
        if !ok {
            t.Errorf("Expected evaluated object to be type of object.Integer. Got %T instead", evaluated)
            continue
        }
        if val.Value != test.expected {
            t.Errorf("Expected function call to return object.Integer with value: %d but found %d instead",
                test.expected, val.Value)
        }
    }
}
