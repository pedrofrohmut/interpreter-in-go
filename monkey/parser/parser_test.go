// monkey/parser/parser_test.go

package parser

import (
    "fmt"
    "testing"
    "bytes"
    _"strconv"
    "monkey/lexer"
    "monkey/ast"
)

func checkParserErrors(t *testing.T, parser *Parser) {
    for i, err := range parser.errors {
        fmt.Printf("# [%d] - ERROR: %s\n", i, err)
    }
    if len(parser.errors) > 0 {
        t.Fatalf("Parser errors")
    }
}

func testIdentifier(t *testing.T, expression ast.Expression, expected string) {
    var iden, ok = expression.(*ast.Identifier)
    if !ok {
        t.Errorf("Expression is not an Identifier. Got %T instead", expression)
        return
    }
    if iden.Value != expected {
        t.Errorf("Expected Identifier value to be %s but found %s instead", expected, iden.Value)
    }
}

func testIntegerLiteral(t *testing.T, expression ast.Expression, expectedValue int64) {
    var intLit, ok = expression.(*ast.IntegerLiteral)
    if !ok {
        t.Errorf("Expression is not an IntegerLiteral. Got %T instead", expression)
        return
    }
    if intLit.Value != expectedValue {
        t.Errorf("Expected IntegerLiteral value to be %d but got %d instead", expectedValue, intLit.Value)
    }
}

func testBooleanLiteral(t *testing.T, expression ast.Expression, expected bool) {
    var boo, ok = expression.(*ast.Boolean)
    if !ok {
        t.Errorf("Expression is not a Boolean. Got %T instead", expression)
        return
    }
    if boo.Value != expected {
        t.Errorf("Expected Boolean value to be '%t' but got '%t' instead", expected, boo.Value)
    }
}

func testLiteralExpression(t *testing.T, expression ast.Expression, expected any) {
    switch tmp := expected.(type) {
    case int:
        testIntegerLiteral(t, expression, int64(tmp))
    case int64:
        testIntegerLiteral(t, expression, tmp)
    case string:
        testIdentifier(t, expression, tmp)
    case bool:
        testBooleanLiteral(t, expression, tmp)
    default:
        t.Errorf("Tested expression type is not a covered on testLiteralExpression. Got %T", expression)
    }
}

func testInfixExpression(t *testing.T, expression ast.Expression, expectedLeft any,
        expectedOperator string, expectedRight any) {
    var inf, ok = expression.(*ast.InfixExpression)
    if !ok {
        t.Errorf("Expression is not an InfixExpression. Got %T instead", expression)
        return
    }
    // Test Left Value
    testLiteralExpression(t, inf.Left, expectedLeft)
    // Test Operator
    if inf.Operator != expectedOperator {
        t.Errorf("Expected Operator to be '%s' but got '%s' instead", expectedOperator, inf.Operator)
    }
    // Test Right Value
    testLiteralExpression(t, inf.Right, expectedRight)
}

func TestLetStatements(t *testing.T) {
    input := `
        let x = 5;
        let y = 10;
        let z = 15;
    `
    lexer := lexer.NewLexer(input)
    parser := NewParser(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)
    if len(program.Statements) != 3 {
        t.Fatalf("Expected program to have %d statements but got %d instead", 3, len(program.Statements))
    }
    program.PrintStatements()
}

func TestReturnStatements(t *testing.T) {
    input := `
        return 5;
        return 10;
        return 15;
    `
    lexer := lexer.NewLexer(input)
    parser := NewParser(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)
    if len(program.Statements) != 3 {
        t.Fatalf("Expected program to have %d statements but got %d instead", 3, len(program.Statements))
    }
    program.PrintStatements()
}

func TestIdentifierExpression(t *testing.T) {
    input := "foobar;"
    lexer := lexer.NewLexer(input)
    parser := NewParser(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)
    if len(program.Statements) != 1 {
        t.Fatalf("Expected program to have %d statements but got %d instead", 1, len(program.Statements))
    }
    program.PrintStatements()

    stm, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Program first statement is not an ExpressionStatement, got %s instead", program.Statements[0])
    }

    ident, ok := stm.Expression.(*ast.Identifier)
    if !ok {
        t.Fatalf("Statement expression is not an Identifier, got %s instead", stm.Expression)
    }

    if ident.Value != "foobar" {
        t.Errorf("Expected identifier value to be '%s' but '%s' instead", "foobar", ident.Value)
    }
}

func TestIntegerExpression(t *testing.T) {
    input := "1234;"
    lexer := lexer.NewLexer(input)
    parser := NewParser(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)
    if len(program.Statements) != 1 {
        t.Fatalf("Expected program to have %d statements but got %d instead", 1, len(program.Statements))
    }
    program.PrintStatements()

    stm, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Program first statement is not an ExpressionStatement, got %s instead", program.Statements[0])
    }

    liter, ok := stm.Expression.(*ast.IntegerLiteral)
    if !ok {
        t.Fatalf("Statement expression is not an IntegerLiteral, got %T instead", stm.Expression)
    }

    if liter.Value != 1234 {
        t.Errorf("Expected integer literal value to be '%d' but got '%d' instead", 1234, liter.Value)
    }
}

func TestParsingPrefixExpression(t *testing.T) {
    tests := []struct {
        input string; operator string; value any
    } {
        { "!5",     "!", 5     },
        { "-15",    "-", 15    },
        { "!true",  "!", true  },
        { "!false", "!", false },
    }
    var acc bytes.Buffer
    for _, x := range tests { acc.WriteString(x.input + ";\n") }
    input := acc.String()
    lexer := lexer.NewLexer(input)
    parser := NewParser(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)
    if len(program.Statements) != len(tests) {
        t.Fatalf("Expected program to have %d statements but got %d instead", len(tests), len(program.Statements))
    }
    program.PrintStatements()

    for i, test := range tests {
        stm, ok := program.Statements[i].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("Program first statement is not an ExpressionStatement, got %s instead", program.Statements[0])
        }

        pref, ok := stm.Expression.(*ast.PrefixExpression)
        if !ok {
            t.Fatalf("Statement expression is not a prefix expression, got %T instead", stm.Expression)
        }

        // Test Prefix Operator
        if pref.Operator != test.operator {
            t.Errorf("Expected prefix expression operator to be %s but got %s instead", test.operator, pref.Operator)
        }

        // Test the Prefix Value
        testLiteralExpression(t, pref.Value, test.value)
    }
}

func TestParsingInfixExpression(t *testing.T) {
    tests := []struct {
        input string; left any; operator string; right any
    } {
        { "5 + 5",          5,     "+",  5     },
        { "5 - 5",          5,     "-",  5     },
        { "5 * 5",          5,     "*",  5     },
        { "5 / 5",          5,     "/",  5     },
        { "5 < 5",          5,     "<",  5     },
        { "5 > 5",          5,     ">",  5     },
        { "5 == 5",         5,     "==", 5     },
        { "5 != 5",         5,     "!=", 5     },

        // Booleans
        { "true == true",   true,  "==", true  },
        { "true != false",  true,  "!=", false },
        { "false == false", false, "==", false },
    }
    var acc bytes.Buffer
    for _, x := range tests { acc.WriteString(x.input + ";\n") }
    input := acc.String()
    lexer := lexer.NewLexer(input)
    parser := NewParser(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)
    if len(program.Statements) != len(tests) {
        t.Fatalf("Expected program to have %d statements but got %d instead", len(tests), len(program.Statements))
    }
    program.PrintStatements()

    for i, test := range tests {
        stm, ok := program.Statements[i].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("Program first statement is not an ExpressionStatement, got %s instead", program.Statements[0])
        }

        testInfixExpression(t, stm.Expression, test.left, test.operator, test.right)
    }
}

func TestOperatorPrecedenceParsing(t *testing.T) {
    tests := []struct {
          input string;                 expected string
    } {
        { "-a * b",                     "((-a) * b)"                             },
        { "!-a",                        "(!(-a))"                                },
        { "a + b + c",                  "((a + b) + c)"                          },
        { "a + b - c",                  "((a + b) - c)"                          },
        { "a * b * c",                  "((a * b) * c)"                          },
        { "a * b / c",                  "((a * b) / c)"                          },
        { "a + b / c",                  "(a + (b / c))"                          },
        { "a + b * c + d / e - f",      "(((a + (b * c)) + (d / e)) - f)"        },
        { "3 + 4",                      "(3 + 4)"                                },
        { "-5 * 5",                     "((-5) * 5)"                             },
        { "5 > 4 == 3 < 4",             "((5 > 4) == (3 < 4))"                   },
        { "5 < 4 != 3 > 4",             "((5 < 4) != (3 > 4))"                   },
        { "3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))" },

        // Booleans
        { "true",                       "true"                                   },
        { "false",                      "false"                                  },
        { "3 > 5 == false",             "((3 > 5) == false)"                     },
        { "3 < 5 == true",              "((3 < 5) == true)"                      },

        // My Custom tests
        { "a + b + c + d",              "(((a + b) + c) + d)"                    },
        { "a + b * c",                  "(a + (b * c))"                          },
        { "a + b * c * d",              "(a + ((b * c) * d))"                    },
        { "a + b * c * d - e",          "((a + ((b * c) * d)) - e)"              },

        // Grouped Expressions
        { "1 + (2 + 3) + 4",            "((1 + (2 + 3)) + 4)"                    },
        { "(5 + 5) * 2",                "((5 + 5) * 2)"                          },
        { "2 / (5 + 5)",                "(2 / (5 + 5))"                          },
        { "-(5 + 5)",                   "(-(5 + 5))"                             },
        { "!(true == true)",            "(!(true == true))"                      },
    }
    var acc bytes.Buffer
    for _, test := range tests { acc.WriteString(test.input + ";\n") }
    input := acc.String()
    lexer := lexer.NewLexer(input)
    parser := NewParser(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)
    if len(program.Statements) != len(tests) {
        t.Fatalf("Expected programs to have %d statements but got %d instead", len(tests), len(program.Statements))
    }
    program.PrintStatements()

    for i, test := range tests {
        stm, ok := program.Statements[i].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("[%d] Program statement is not a ExpressionStatement, got %s instead", i, program.Statements[i])
        }

        var stmStr = stm.String()
        var expected = test.expected + ";"
        if stmStr != expected {
            t.Errorf("[%d] Expected statement to be '%s', but got '%s' instead", i, stmStr, expected)
        }
    }
}

func TestParsingIfExpressions(t *testing.T) {
    var input = "if (x < y) { x };"
    var lexer = lexer.NewLexer(input)
    var parser = NewParser(lexer)
    var program = parser.ParseProgram()

    checkParserErrors(t, parser)
    if len(program.Statements) != 1 {
        t.Fatalf("Expected programs to have %d statements but got %d instead", 1, len(program.Statements))
    }
    program.PrintStatements()
}

func TestParsingIfExpressions2(t *testing.T) {
    var input = "if (x < y) { x } else { y };"
    var lexer = lexer.NewLexer(input)
    var parser = NewParser(lexer)
    var program = parser.ParseProgram()

    checkParserErrors(t, parser)
    if len(program.Statements) != 1 {
        t.Fatalf("Expected programs to have %d statements but got %d instead", 1, len(program.Statements))
    }
    program.PrintStatements()
}
