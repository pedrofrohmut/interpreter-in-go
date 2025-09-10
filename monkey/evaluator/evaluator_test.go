// monkey/evaluator/evaluator_test.go

package evaluator

import (
    "monkey/lexer"
    "monkey/object"
    "monkey/parser"
    "monkey/ast"
    "monkey/test_utils"
    "testing"
)

func getParsedProgram(input string) (*ast.Program, *parser.Parser) {
    var lexer = lexer.NewLexer(input)
    var parser = parser.NewParser(lexer)
    var program = parser.ParseProgram()
    return program, parser
}

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

        // Strings
        { `"Hello, " - "World!"`, "unknown operator: String - String" },
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
        { // Using closures (a and b from outer scope)
            `
                let a = 5;
                let b = 10;
                let add = fn(x) { a + b + x; };
                add(10);
            `,
            25,
        },
        { // Using identifiers on the callExpression instead of integer literals
            `
                let a = 5;
                let b = 10;
                let add = fn(x, y) { return x + y; };
                add(a, b);
            `,
            15,
        },
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

func TestClosures(t *testing.T) {
    var input = `
        let newAdder = fn (x) {
            return fn (y) { return x + y; };
        };
        let addTwo = newAdder(2);
        addTwo(5);
    `
    var expected int64 = 7
    var lexer = lexer.NewLexer(input)
    var parser = parser.NewParser(lexer)
    var program = parser.ParseProgram()

    test_utils.CheckForParserErrors(t, parser)

    var env = object.NewEnvironment()
    var evaluated = Eval(program, env)
    if test_utils.CheckForEvalError(t, evaluated) { return }

    var val, ok = evaluated.(*object.Integer)
    if !ok {
        t.Errorf("Expected evaluated object to be type of object.Integer. Got %T instead", evaluated)
        return
    }
    if val.Value != expected {
        t.Errorf("Expected evaluated value object value to be %d but got %d instead", expected, val.Value)
    }
}

func TestHelloWorld(t *testing.T) {
    var input = `"Hello, World!"`
    var lexer = lexer.NewLexer(input)
    var parser = parser.NewParser(lexer)
    var program = parser.ParseProgram()

    test_utils.CheckForParserErrors(t, parser)

    var env = object.NewEnvironment()
    var evaluated = Eval(program, env)
    if test_utils.CheckForEvalError(t, evaluated) { return }

    var val, ok = evaluated.(*object.String)
    if !ok {
        t.Errorf("Expected evaluated object to be type of object.String. Got %T instead", evaluated)
        return
    }
    var expected = "Hello, World!"
    if val.Value != expected {
        t.Errorf("Expected evaluated value object value to be '%s' but got '%s' instead", expected, val.Value)
    }
}

func TestStringConcat(t *testing.T) {
    var input = `"Hello, " + "World!";`
    var lexer = lexer.NewLexer(input)
    var parser = parser.NewParser(lexer)
    var program = parser.ParseProgram()

    test_utils.CheckForParserErrors(t, parser)

    var env = object.NewEnvironment()
    var evaluated = Eval(program, env)
    if test_utils.CheckForEvalError(t, evaluated) { return }

    var val, ok = evaluated.(*object.String)
    if !ok {
        t.Errorf("Expected evaluated object to be type of object.String. Got %T instead", evaluated)
        return
    }
    var expected = "Hello, World!"
    if val.Value != expected {
        t.Errorf("Expected evaluated value object value to be '%s' but got '%s' instead", expected, val.Value)
    }
}

func TestBuiltinFunctionsLen(t *testing.T) {
    var tests = []struct {
        input string; expected any
    } {
        { `len("")`,                    0                                                 },
        { `len("four")`,                4                                                 },
        { `len("hello, world!")`,       13                                                },
        { `len(1)`,                     "argument to len not supported, got Integer"      },
        { `len("one", "two")`,          "wrong number of arguments. expected=1 but got=2" },
        { `len("one", "two", "three")`, "wrong number of arguments. expected=1 but got=3" },

        // Arrays
        { `len([])`,                      0 },
        { `len(["one"])`,                 1 },
        { `len(["one", "two"])`,          2 },
        { `len(["one", "two", "three"])`, 3 },
    }


    for _, test := range tests {
        var lexer = lexer.NewLexer(test.input)
        var parser = parser.NewParser(lexer)
        var program = parser.ParseProgram()
        test_utils.CheckForParserErrors(t, parser)

        var evaluated = Eval(program, object.NewEnvironment())
        // if test_utils.CheckForEvalError(t, evaluated) { continue }

        switch expected := test.expected.(type) {
        case int:
            var integer, ok = evaluated.(*object.Integer)
            if !ok {
                t.Errorf("Expected evaluated to be object.Integer but got %T instead", evaluated)
                continue
            }
            if integer.Value != int64(expected) {
                t.Errorf("Expected evaluated object.Integer value to be %d but got %d instead",
                    test.expected, integer.Value)
            }
        case string:
            var strObj, ok = evaluated.(*object.Error)
            if !ok {
                t.Errorf("Expected evaluated to be object.String but got %T instead", evaluated)
                continue
            }
            if strObj.Message != expected {
                t.Errorf("Expected evaluated object.String value to be '%s' but got '%s' instead", expected, strObj.Message)
            }
        }
    }
}

func TestBuiltinFirst(t *testing.T) {
    var tests = []struct {
        input string; expected any
    } {
        { `first("")`,           nil },
        { `first("foo")`,        'f' },
        { `first([])`,           nil },
        { `first([61])`,         61  },
        { `first([61, 62])`,     61  },
        { `first([61, 62, 63])`, 61  },
    }

    for _, test := range tests {
        var lexer = lexer.NewLexer(test.input)
        var parser = parser.NewParser(lexer)
        var program = parser.ParseProgram()
        test_utils.CheckForParserErrors(t, parser)

        var evaluated = Eval(program, object.NewEnvironment())
        if test_utils.CheckForEvalError(t, evaluated) { continue }

        switch expected := test.expected.(type) {
        case int:
            var intObj, ok = evaluated.(*object.Integer)
            if !ok {
                t.Errorf("Expected evaluated to be type object.Integer but got %T instead", evaluated)
                continue
            }
            if intObj.Value != int64(expected) {
                t.Errorf("Expected object.Integer value to be %d but got %d instead", expected, intObj.Value)
            }
        case int32:
            var ch, ok = evaluated.(*object.Char)
            if !ok {
                t.Errorf("Expected evaluated to be type object.Char but got %T instead", evaluated)
                continue
            }
            if ch.Value != byte(expected) {
                t.Errorf("Expected object.Char value to be %d but got %d instead", expected, ch.Value)
            }
        case nil:
            var _, ok = evaluated.(*object.Null)
            if !ok {
                t.Errorf("Expected evaluated to be of type object.Null but got %T instead", evaluated)
            }
        default:
            t.Errorf("Type %T not covered", test.expected)
        }
    }
}

func TestBuiltinLast(t *testing.T) {
        var tests = []struct {
        input string; expected any
    } {
        { `last("")`,           nil },
        { `last("foobar")`,     'r' },
        { `last([])`,           nil },
        { `last([61])`,         61  },
        { `last([61, 62])`,     62  },
        { `last([61, 62, 63])`, 63  },
    }

    for _, test := range tests {
        var lexer = lexer.NewLexer(test.input)
        var parser = parser.NewParser(lexer)
        var program = parser.ParseProgram()
        test_utils.CheckForParserErrors(t, parser)

        var evaluated = Eval(program, object.NewEnvironment())
        if test_utils.CheckForEvalError(t, evaluated) { continue }

        switch expected := test.expected.(type) {
        case int:
            var intObj, ok = evaluated.(*object.Integer)
            if !ok {
                t.Errorf("Expected evaluated to be type object.Integer but got %T instead", evaluated)
                continue
            }
            if intObj.Value != int64(expected) {
                t.Errorf("Expected object.Integer value to be %d but got %d instead", expected, intObj.Value)
            }
        case int32:
            var ch, ok = evaluated.(*object.Char)
            if !ok {
                t.Errorf("Expected evaluated to be type object.Char but got %T instead", evaluated)
                continue
            }
            if ch.Value != byte(expected) {
                t.Errorf("Expected object.Char value to be %d but got %d instead", expected, ch.Value)
            }
        case nil:
            var _, ok = evaluated.(*object.Null)
            if !ok {
                t.Errorf("Expected evaluated to be of type object.Null but got %T instead", evaluated)
            }
        default:
            t.Errorf("Type %T not covered", test.expected)
        }
    }
}

func TestBuiltinRest(t *testing.T) {
    var tests = []struct {
        input string; expected any
    } {
        { `rest("")`,           ""             },
        { `rest("foobar")`,     "oobar"        },
        { `rest([])`,           []int {}       },
        { `rest([61])`,         []int {}       },
        { `rest([61, 62])`,     []int {62}     },
        { `rest([61, 62, 63])`, []int {62, 63} },
    }

    for _, test := range tests {
        var lexer = lexer.NewLexer(test.input)
        var parser = parser.NewParser(lexer)
        var program = parser.ParseProgram()
        test_utils.CheckForParserErrors(t, parser)

        var evaluated = Eval(program, object.NewEnvironment())
        if test_utils.CheckForEvalError(t, evaluated) { continue }

        switch expected := test.expected.(type) {
        case []int:
            var arrObj, ok = evaluated.(*object.Array)
            if !ok {
                t.Errorf("Expected evaluated to be of type object.Array but got %T instead", evaluated)
                continue
            }
            for i, obj := range arrObj.Elements {
                var intObj, okInt = obj.(*object.Integer)
                if !okInt {
                    t.Errorf("Expected array element to be object.Integer but got %T instead", obj)
                    continue
                }
                if intObj.Value != int64(expected[i]) {
                    t.Errorf("Expected array element value to be %d but got %d instead", expected[i], intObj.Value)
                }
            }
        case string:
            var strObj, ok = evaluated.(*object.String)
            if !ok {
                t.Errorf("Expected evaluated to be of type object.String but got %T instead", evaluated)
                continue
            }
            if strObj.Value != expected {
                t.Errorf("Expected object string value to be '%s' but got '%s' instead", expected, strObj.Value)
            }
        default:
            t.Errorf("Type %T not covered", test.expected)
        }
    }
}

func TestBuiltinPush(t *testing.T) {
    var inputs = []string {
        `let myarr = [1, 2, 3, 4, 5]; push(myarr, 666); myarr;`,
        `let myarr = push([1, 2, 3, 4, 5], 666); myarr;`,
    }
    for _, input := range inputs {
        var lexer = lexer.NewLexer(input)
        var parser = parser.NewParser(lexer)
        var program = parser.ParseProgram()
        test_utils.CheckForParserErrors(t, parser)

        var evaluated = Eval(program, object.NewEnvironment())
        if test_utils.CheckForEvalError(t, evaluated) { return }

        var arr, okArr = evaluated.(*object.Array)
        if !okArr {
            t.Errorf("Expected evaluated to be of type object.Array but got %T instead", evaluated)
            return
        }

        var expectations =  []int64 {1, 2, 3, 4, 5, 666}
        for i, expected := range expectations {
            var elem = arr.Elements[i].(*object.Integer)
            if elem.Value != expected {
                t.Errorf("Expected array elements[%d] to be %d but got %d instead", i, expected, elem.Value)
            }
        }
    }
}

func TestArrayLiterals(t *testing.T) {
    var input = `[1, 2 * 3, 4 + 5]`

    var lexer = lexer.NewLexer(input)
    var parser = parser.NewParser(lexer)
    var program = parser.ParseProgram()
    test_utils.CheckForParserErrors(t, parser)

    var evaluated = Eval(program, object.NewEnvironment())
    if test_utils.CheckForEvalError(t, evaluated) { return }

    var expectations = []int64 {1, 6, 9}
    for i, expected := range expectations {
        var arr, okArr = evaluated.(*object.Array)
        if !okArr {
            t.Errorf("Expected evaluated to be type of object.Array but got %T instead", evaluated)
            continue
        }
        var check = arr.Elements[i].(*object.Integer)
        if check.Value != expected {
            t.Errorf("Expected array object.Integer value to be %d but got %d instead", expected, check.Value)
        }
    }
}

func TestIndexExpression(t *testing.T) {
    var tests = []struct {
        input string; expected any
    } {
        { "[1, 2, 3][0];",                                         1   },
        { "[1, 2, 3][1];",                                         2   },
        { "[1, 2, 3][2];",                                         3   },
        { "let i = 0; [1][i];",                                    1   },
        { "[1, 2, 3][1 + 1];",                                     3   },
        { "let myArr = [1, 2, 3]; myArr[2];",                      3   },
        { "let myArr = [1, 2, 3]; myArr[0] + myArr[1] + myArr[2]", 6   },
        { "let myArr = [1, 2, 3]; let i = myArr[0]; myArr[i];",    2   },

        // Out of bounds check
        { "[1, 2, 3][3]",                                          nil },
        { "[1, 2, 3][-1]",                                         nil },
        { "let myArr = [1, 2, 3]; myArr[3];",                      nil },
        { "let myArr = [1, 2, 3]; myArr[-1];",                     nil },
    }

    for _, test := range tests {
        var lexer = lexer.NewLexer(test.input)
        var parser = parser.NewParser(lexer)
        var program = parser.ParseProgram()
        if test_utils.CheckForParserErrors(t, parser) { continue }
        // program.PrintStatements()

        var evaluated = Eval(program, object.NewEnvironment())
        if test_utils.CheckForEvalError(t, evaluated) { continue }

        switch x := test.expected.(type) {
        case int:
            var intObj, ok = evaluated.(*object.Integer)
            if !ok {
                t.Errorf("Expected evaluated to be of type object.Integer but got %T instead", evaluated)
                continue
            }
            if intObj.Value != int64(x) {
                t.Errorf("Expected evaluated value to be %d but got %d instead", int64(x), intObj.Value)
            }
        case nil:
            if evaluated.Type() != object.NullType {
                t.Errorf("Expected evaluated object to be object.NullType but got %s instead", evaluated.Type())
            }
        }
    }
}

func TestHashLiterals(t *testing.T) {
    var input = `
        let two = "two";
        {
            "one": 10 - 9,
            two: 1 + 1,
            "thr" + "ee": 6 / 2,
            4: 4,
            true: 5,
            false: 6
        };
    `
    var lexer = lexer.NewLexer(input)
    var parser = parser.NewParser(lexer)
    var program = parser.ParseProgram()
    test_utils.CheckForParserErrors(t, parser)

    var evaluated = Eval(program, object.NewEnvironment())
    if test_utils.CheckForEvalError(t, evaluated) { return }

    var hash, okHash = evaluated.(*object.Hash)
    if !okHash {
        t.Errorf("Expected evaluated to be object of type Hash but got %T instead", evaluated)
        return
    }

    var expectations = map[uint64]int64 {
        (&object.String { Value: "one" }).HashKey().Value:   1,
        (&object.String { Value: "two" }).HashKey().Value:   2,
        (&object.String { Value: "three" }).HashKey().Value: 3,
        (&object.Integer { Value: 4 }).HashKey().Value:      4,
        ObjTrue.HashKey().Value:                             5,
        ObjFalse.HashKey().Value:                            6,
    }

    for key, value := range hash.Pairs {
        var expectedValue, ok = expectations[key.Value]
        if !ok {
            t.Errorf("Expected value to be found in expectations to be found with current key.value but nothing was found")
            continue
        }

        var objInt, okObjInt = value.Value.(*object.Integer)
        if !okObjInt {
            t.Errorf("Expected pair value to be of type object.Integer but got %T instead", value)
            continue
        }
        if expectedValue != objInt.Value {
            t.Errorf("Expected pair value to be %d but got %d instead", expectedValue, objInt.Value)
        }
    }
}

func TestHashIndexExpression(t *testing.T) {
    var tests = []struct {
        input string; expected any
    }{
        { `{ "foo": 5 }["foo"]`,                5   },
        { `{ "foo": 5 }["bar"}`,                nil },
        { `let key = "foo"; { "foo": 5 }[key]`, 5   },
        { `{}["foo"]`,                          nil },
        { `{ 5: 5 }[5]`,                        5   },
        { `{ true: 5 }[true]`,                  5   },
        { `{ false: 5 }[false]`,                5   },
    }

    for _, test := range tests {
        var program, parser = getParsedProgram(test.input)
        test_utils.CheckForParserErrors(t, parser)

        var evaluated = Eval(program, object.NewEnvironment())
        if test_utils.CheckForEvalError(t, evaluated) { continue }

        switch x := test.expected.(type) {
        case int:
            var objInt, ok = evaluated.(*object.Integer)
            if !ok {
                t.Errorf("Expected evaluated be of type object.Integer but found %T instead", evaluated)
                continue
            }
            if objInt.Value != int64(x) {
                t.Errorf("Expected evaluated value to be %d but gut %d instead", x, objInt.Value)
            }
        case nil:
            var _, ok = evaluated.(*object.Null)
            if !ok {
                t.Errorf("Expected evaluated to be of type object.Null but got %T instead", evaluated)
            }
        }
    }
}

func TestBuiltinPuts(t *testing.T) {
    var input = `puts("Hello, World!")`
    var program, parser = getParsedProgram(input)
    test_utils.CheckForParserErrors(t, parser)

    var evaluated = Eval(program, object.NewEnvironment())
    if test_utils.CheckForEvalError(t, evaluated) { return }

    _ = program
}
