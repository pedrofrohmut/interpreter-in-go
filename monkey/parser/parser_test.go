// monkey/parser/parser_test.go

package parser

import (
    "fmt"
    "testing"
    "bytes"
    "strconv"
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
        input string; operator string; value int64
    } {
        { "!5",  "!", 5  },
        { "-15", "-", 15 },
    }
    var acc bytes.Buffer
    for _, x := range tests { acc.WriteString(x.input + ";\n") }
    input := acc.String()
    lexer := lexer.NewLexer(input)
    parser := NewParser(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)
    if len(program.Statements) != 2 {
        t.Fatalf("Expected program to have %d statements but got %d instead", 2, len(program.Statements))
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

        if pref.Value.String() != strconv.FormatInt(test.value, 10) {
            t.Errorf("Expected prefix expression value to be %d but got %d instead", test.value, pref.Value)
        }

        if pref.Operator != test.operator {
            t.Errorf("Expected prefix expression operator to be %s but got %s instead", test.operator, pref.Operator)
        }
    }
}

func TestParsingInfixExpression(t *testing.T) {
    tests := []struct {
        input string; left int64; operator string; right int64
    } {
        { "5 + 5",  5, "+",  5 },
        { "5 - 5",  5, "-",  5 },
        { "5 * 5",  5, "*",  5 },
        { "5 / 5",  5, "/",  5 },
        { "5 < 5",  5, "<",  5 },
        { "5 > 5",  5, ">",  5 },
        { "5 == 5", 5, "==", 5 },
        { "5 != 5", 5, "!=", 5 },
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

        infix, ok := stm.Expression.(*ast.InfixExpression)
        if !ok {
            t.Fatalf("Statement expression is not a infix expression, got %T instead", stm.Expression)
        }

        leftValue := infix.Left.(*ast.IntegerLiteral).Value
        if leftValue != test.left {
            t.Errorf("Expected infix expression left value to be %d but got %d instead", test.left, leftValue)
        }

        if infix.Operator != test.operator {
            t.Errorf("Expected infix expression operator to be %s but got %s instead", test.operator, infix.Operator)
        }

        rightValue := infix.Right.(*ast.IntegerLiteral).Value
        if rightValue != test.right {
            t.Errorf("Expected infix expression right value to be %d but got %d instead", test.right, rightValue)
        }
    }
}

func TestOperatorPrecedence(t *testing.T) {
    tests := []struct {
          input string;                 expected string
    } {
        // { "!-a",                        "(!(-a))"                         },
        { "-a * b",                     "((-a) * b)"                      },
        { "a + b + c",                  "((a + b) + c)"                      },
        { "a + b - c",                  "((a + b) - c)"                   },
        { "a * b * c",                  "((a * b) * c)"                   },
        { "a * b / c",                  "((a * b) / c)"                   },
        { "a + b / c",                  "(a + (b / c))"                   },
        { "a + b * c + d / e - f",      "(((a + (b * c)) + (d / e)) - f)" },
        { "3 + 4",                      "(3 + 4)"                         },
        { "-5 * 5",                     "((-5) * 5)"                      },

        // My Custom tests
        { "a + b + c + d",              "(((a + b) + c) + d)"                },
        { "a + b * c",                  "(a + (b * c))"                      },
        { "a + b * c * d",              "(a + ((b * c) * d))"                },
        { "a + b * c * d - e",          "((a + ((b * c) * d)) - e)"          },
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
}
