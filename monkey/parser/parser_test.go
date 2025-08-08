// monkey/parser/parser_test.go

package parser

import (
    "bytes"
    "fmt"
    "strconv"
    "testing"
    "monkey/ast"
    "monkey/lexer"
    "monkey/token"
)

func checkParserErrors(t *testing.T, par *Parser) {
    errors := par.Errors()
    if len(errors) == 0 { return }
    for i, err := range errors {
        t.Errorf("[%d] Parser Error: %s\n", i, err)
    }
    t.FailNow()
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, expectedValue int64) {
    intLiteral, ok := exp.(*ast.IntegerLiteral)
    if !ok {
        t.Errorf("Expression is not an Integer Literal")
        return
    }
    if intLiteral.Value != expectedValue {
        t.Errorf("Expected Integer Literal value to be %d but got %d instead", expectedValue, intLiteral.Value)
    }
    expectedStr := strconv.FormatInt(expectedValue, 10)
    if intLiteral.TokenLiteral() != expectedStr {
        t.Errorf("Expected Integer token literal to be '%s' but got '%s' instead", expectedStr, intLiteral.TokenLiteral())
    }
}

func testIdentifier(t *testing.T, exp ast.Expression, expectedValue string) {
    ident, ok := exp.(*ast.Identifier)
    if !ok {
        t.Errorf("Expression is not an Identifier")
        return
    }
    if ident.Value != expectedValue {
        t.Errorf("Expected Identifier Value to be %s but got %s instead", expectedValue, ident.Value)
    }
    if ident.TokenLiteral() != expectedValue {
        t.Errorf("Expected Identifier TokenLiteral to be %s but got %s instead", expectedValue, ident.TokenLiteral())
    }
}

func testBoolean(t *testing.T, exp ast.Expression, expectedValue bool) {
    boolean, ok := exp.(*ast.Boolean)
    if !ok {
        t.Errorf("Expression is not a Boolean")
        return
    }
    if boolean.Value != expectedValue {
        t.Errorf("Expected Boolean value to be %t but got %t instead", boolean.Value, expectedValue)
    }
    strValue := fmt.Sprintf("%t", expectedValue)
    if boolean.TokenLiteral() != strValue {
        t.Errorf("Expected Boolean TokenLiteral to be %s but got %s instead", boolean.TokenLiteral(), strValue)
    }
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) {
    switch v := expected.(type) {
    case int:
        testIntegerLiteral(t, exp, int64(v))
    case int64:
        testIntegerLiteral(t, exp, v)
    case bool:
        testBoolean(t, exp, v)
    case string:
        testIdentifier(t, exp, v)
    default:
        t.Errorf("Type of exp not handled. got %T", exp)
    }
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any, expectedOperator string, right any) {
    infExp, ok := exp.(*ast.InfixExpression)
    if !ok {
        t.Errorf("Expression is not an Infix Expression")
        return
    }
    testLiteralExpression(t, infExp.Left, left)
    if infExp.Operator != expectedOperator {
        t.Errorf("Expected operator to be %s but got %s instead", expectedOperator, infExp.Operator)
    }
    testLiteralExpression(t, infExp.Right, right)
}

func TestLetStatement(t *testing.T) {
    input := `
        let x = 5;
        let y = 10;
        let z = 15;
    `
    lex := lexer.NewLexer(input)
    par := NewParser(lex)
    pro := par.ParseProgram()

    pro.PrintStatements()

    checkParserErrors(t, par)
    if pro == nil {
        t.Fatalf("Program is nil")
    }
    if len(pro.Statements) != 3 {
        t.Fatalf("Expected program to have %d statements but got %d\n", 3, len(pro.Statements))
    }

    tests := []struct { expectedIdentifier string } {
        {"x"}, {"y"}, {"z"},
    }

    for i, test := range tests {
        if pro.Statements[i].TokenLiteral() != "let" {
            t.Errorf("Expected statement literal to be '%s' but got '%s' instead",
                "let", pro.Statements[i].TokenLiteral())
        }
        stm, ok := pro.Statements[i].(*ast.LetStatement)
        if !ok {
            t.Errorf("Is not a LetStatement")
        }
        if stm.Token.Type != token.LET {
            t.Errorf("[%dl] Expected identifier to be %s but got %s", i, token.LET, stm.Token.Type)
        }
        if stm.Identifier.Value != test.expectedIdentifier {
            t.Errorf("[%d] Expected identifier to be %s but got %s", i, stm.Identifier.Value, test.expectedIdentifier)
        }
    }
}

func TestReturnStatement(t *testing.T) {
    tests := []struct {
          input string;            expected string
    }{
        { "return;",               "" },
        { "return 5;",             "5" },
        { "return 10;",            "10" },
        { "return 1234;",          "1234" },
        { "return 1 + 2 / 3 * 4;", "(1 + ((2 / 3) * 4))" },
        { "return 1 + 2 + 3 + 4;", "(((1 + 2) + 3) + 4)" },
    }
    var acc bytes.Buffer
    for _, x := range tests {
        acc.WriteString(x.input + "\n")
    }
    input := acc.String()
    lex := lexer.NewLexer(input)
    par := NewParser(lex)
    pro := par.ParseProgram()

    pro.PrintStatements()

    checkParserErrors(t, par)
    if pro == nil {
        t.Fatalf("Program is nil")
    }
    if len(pro.Statements) != 6 {
        t.Fatalf("Expected program to have %d statements but got %d\n", 6, len(pro.Statements))
    }

    for i, test := range tests {
        stm, ok := pro.Statements[i].(*ast.ReturnStatement)
        if !ok {
            t.Errorf("Is not a ReturnStatement")
        }

        if stm.Token.Type != token.RETURN {
            t.Errorf("[%d] Expected token type to be %s but got %s", i, token.RETURN, stm.Token.Type)
        }
        if stm.Token.Literal != "return" {
            t.Errorf("[%d] Expected token literal to be '%s' but got '%s'", i, "return", stm.Token.Literal)
        }
        if stm.TokenLiteral() != "return" {
            t.Errorf("[%d] Expected TokenLiteral to be '%s' but got '%s'", i, "return", stm.Token.Literal)
        }

        if stm.Expression != nil && stm.Expression.String() != test.expected {
            t.Errorf("[%d] Expected return expression to be '%s' but got '%s' instead",
                i, test.expected, stm.Expression.String())
        }
    }
}

func TestErrorsOnLetStatement(t *testing.T) {
    input := `
        let x 5;
        let = 10;
        let 15;
    `
    lex := lexer.NewLexer(input)
    par := NewParser(lex)
    pro := par.ParseProgram()

    if pro == nil {
        t.Fatalf("Program is nil")
    }
    if len(par.errors) != 4 {
        t.Fatalf("Expected number of errors to be %d but got %d instead.", 4, len(par.errors))
    }
}

func TestIdentifierExpression(t *testing.T) {
    input := "foo;"
    lex := lexer.NewLexer(input)
    par := NewParser(lex)
    pro := par.ParseProgram()

    checkParserErrors(t, par)
    if len(pro.Statements) != 1 {
        t.Fatalf("Expected program statements length to be %d but got %d", 1, len(pro.Statements))
    }
    stm, ok := pro.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Not an expression statement")
    }
    ident, ok := stm.Expression.(*ast.Identifier)
    if !ok {
        t.Fatalf("Statement expression is not an identifier")
    }
    if ident.Value != "foo" {
        t.Errorf("Expected identifier value to be '%s' but got '%s' instead", "foo", ident.Value)
    }
    if ident.TokenLiteral() != "foo" {
        t.Errorf("Expected statement expression to be '%s' but got '%s' instead", "foo", ident.TokenLiteral())
    }
}

func TestIntegerLiteralExpression(t *testing.T) {
    input := "5;"
    lex := lexer.NewLexer(input)
    par := NewParser(lex)
    pro := par.ParseProgram()

    if len(pro.Statements) != 1 {
        t.Fatalf("Expected program number of statements to be %d but got %d instead", 1, len(pro.Statements))
    }
    stm, ok := pro.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Expected an ast.ExpressionStatement but got %T instead", pro.Statements[0])
    }
    intLit, ok := stm.Expression.(*ast.IntegerLiteral)
    if !ok {
        t.Fatalf("Expected an ast.IntegerLiteral but got %T instead", stm.Expression)
    }
    if intLit.TokenLiteral() != "5" {
        t.Errorf("Expected statement expression to be '%s' but got '%s' instead", "5", intLit.TokenLiteral())
    }
    if intLit.Token.Type != token.INT {
        t.Errorf("Expected integer literal token type to be %s but got %s instead", token.INT, intLit.Token.Type)
    }
    if intLit.Value != 5 {
        t.Errorf("Expected integer literal value to be %d but got %d instead", 5, intLit.Value)
    }
}

func TestParsingPrefixExpression(t *testing.T) {
    tests := []struct {
        input string; operator string; value any
    } {
        { "!5",     "!", 5 },
        { "-15",    "-", 15 },
        { "!true",  "!", true },
        { "!false", "!", false },
    }
    var acc bytes.Buffer
    for _, test := range tests {
        acc.WriteString(test.input + ";\n")
    }
    input := acc.String()
    lex := lexer.NewLexer(input)
    par := NewParser(lex)
    pro := par.ParseProgram()

    if pro == nil {
        t.Fatalf("Program is nil")
    }
    if len(pro.Statements) != len(tests) {
        t.Fatalf("Expected program number of statements to be %d but got %d instead", len(tests), len(pro.Statements))
    }
    checkParserErrors(t, par)

    for i, test := range tests {
        stm, ok := pro.Statements[i].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("Not an expression statement")
        }
        exp, ok := stm.Expression.(*ast.PrefixExpression)
        if !ok {
            t.Fatalf("Expression statement is not a prefix statement")
        }
        if exp.Operator != test.operator {
            t.Errorf("Expected operator to be '%s' but got '%s' instead", test.operator, exp.Operator)
        }
        // Test Right
        testLiteralExpression(t, exp.Right, test.value)
    }
}

func TestParsingInfixExpression(t *testing.T) {
    tests := []struct {
        input string; leftValue any; operator string; rightValue any
    } {
        { "5 + 5",          5,     "+",  5 },
        { "5 - 5",          5,     "-",  5 },
        { "5 * 5",          5,     "*",  5 },
        { "5 / 5",          5,     "/",  5 },
        { "5 > 5",          5,     ">",  5 },
        { "5 < 5",          5,     "<",  5 },
        { "5 == 5",         5,     "==", 5 },
        { "5 != 5",         5,     "!=", 5 },
        { "true == true",   true,  "==", true },
        { "false == false", false, "==", false },
        { "true != false",  true,  "!=", false },
    }
    var acc bytes.Buffer
    for _, test := range tests {
        acc.WriteString(test.input + ";\n")
    }
    input := acc.String()
    lex := lexer.NewLexer(input)
    par := NewParser(lex)
    pro := par.ParseProgram()

    checkParserErrors(t, par)
    if pro == nil {
        t.Fatalf("Program is nil")
    }
    if len(pro.Statements) != len(tests) {
        t.Fatalf("Expected program number of statements to be %d but got %d instead", len(tests), len(pro.Statements))
    }

    for i, test := range tests {
        stm, ok := pro.Statements[i].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("[%d] Not an expression statement", i)
        }
        expression, ok := stm.Expression.(*ast.InfixExpression)
        if !ok {
            t.Fatalf("[%d] Expression statement is not a infix expression", i)
        }
        // Test Left
        testLiteralExpression(t, expression.Left, test.leftValue)
        // Operator
        if expression.Operator != test.operator {
            t.Fatalf("[%d] Expected Infix Expression Operator to be '%s' but got '%s' instead",
                i, test.operator, expression.Operator)
        }
        // Test Right
        testLiteralExpression(t, expression.Right, test.rightValue)
    }
}

func TestOperatorPrecedenceParsing(t *testing.T) {
    tests := []struct {
          input string;                 expected string
    } {
        { "-a * b",                     "((-a) * b)" },
        { "!-a",                        "(!(-a))" },
        { "a + b + c",                  "((a + b) + c)" },
        { "a + b - c",                  "((a + b) - c)" },
        { "a * b * c",                  "((a * b) * c)" },
        { "a * b / c",                  "((a * b) / c)" },
        { "a + b / c",                  "(a + (b / c))" },
        { "a + b * c + d / e - f",      "(((a + (b * c)) + (d / e)) - f)" },
        { "3 + 4",                      "(3 + 4)" },
        { "-5 * 5",                     "((-5) * 5)" },

        // EQ NOT_EQ
        { "5 > 4 == 3 < 4",             "((5 > 4) == (3 < 4))" },
        { "5 < 4 != 3 > 4",             "((5 < 4) != (3 > 4))" },
        { "3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))" },

        // Booleans
        { "true",                       "true" },
        { "false",                      "false" },
        { "3 > 5 == false",             "((3 > 5) == false)" },
        { "3 < 5 == true",              "((3 < 5) == true)" },

        // Groups
        { "1 + (2 + 3) + 4",            "((1 + (2 + 3)) + 4)" },
        { "(5 + 5) * 2",                "((5 + 5) * 2)" },
        { "2 / (5 + 5)",                "(2 / (5 + 5))" },
        { "-(5 + 5)",                   "(-(5 + 5))" },
        { "!(true == true)",            "(!(true == true))" },

        // Call Expression
        {
            "a + add(b * c) + d",
            "((a + add((b * c))) + d)",
        },
        // {
        //     "add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
        //     "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
        // },
        // {
        //     "add(a + b + c * d / f + g)",
        //     "add((((a + b) + ((c * d) / f)) + g))",
        // },
    }

    var acc bytes.Buffer
    for _, test := range tests {
        acc.WriteString(test.input + ";\n")
    }
    input := acc.String()
    lex := lexer.NewLexer(input)
    par := NewParser(lex)
    pro := par.ParseProgram()

    if pro == nil {
        t.Fatalf("The program is nil")
    }
    if len(pro.Statements) != len(tests) {
        t.Errorf("Expected program to have %d statements but it have %d instead", len(tests), len(pro.Statements))
    }
    checkParserErrors(t, par)

    for i, test := range tests {
        progStr := pro.Statements[i].String()
        // fmt.Println(progStr)
        if progStr != test.expected {
            t.Errorf("Expected program string to be '%s' but got '%s' instead\n", test.expected, progStr)
        }
    }
}

func TestIfExpression(t *testing.T) {
    input := "if (x < y) { x }"
    lex := lexer.NewLexer(input)
    par := NewParser(lex)
    pro := par.ParseProgram()

    if pro == nil {
        t.Fatalf("The program is nil")
    }
    if len(pro.Statements) != 1 {
        t.Errorf("Expected program to have %d statements but it have %d instead", 1, len(pro.Statements))
    }
    checkParserErrors(t, par)

    stm, ok := pro.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("The statement is not an ExpressionStatement. Got %T instead", pro.Statements[0])
    }
    exp, ok := stm.Expression.(*ast.IfExpression)
    if !ok {
        t.Fatalf("The expression is not an IfExpression. Got %T instead", stm.Expression)
    }

    // Check condition
    testInfixExpression(t, exp.Condition, "x", "<", "y")

    // Check consequence
    if len(exp.Consequence.Statements) != 1 {
        t.Errorf("Expected if consequence number of statements to be %d but got %d instead",
            1, len(exp.Consequence.Statements))
    }
    consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Consequence Statement is not an ExpressionStatement. Got %T instead", exp.Consequence.Statements[0])
    }
    testIdentifier(t, consequence.Expression, "x")

    // Check alternative
    if exp.Alternative != nil {
        t.Errorf("Alternative is not nil")
    }
}

func TestIfExpression2(t *testing.T) {
    input := "if (x < y) { x } else { y }"
    lex := lexer.NewLexer(input)
    par := NewParser(lex)
    pro := par.ParseProgram()

    if pro == nil {
        t.Fatalf("The program is nil")
    }
    if len(pro.Statements) != 1 {
        t.Errorf("Expected program to have %d statements but it have %d instead", 1, len(pro.Statements))
    }
    checkParserErrors(t, par)

    stm, ok := pro.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("The statement is not an ExpressionStatement. Got %T instead", pro.Statements[0])
    }
    exp, ok := stm.Expression.(*ast.IfExpression)
    if !ok {
        t.Fatalf("The expression is not an IfExpression. Got %T instead", stm.Expression)
    }

    // Check Condition
    testInfixExpression(t, exp.Condition, "x", "<", "y")

    // Check Consequence
    if len(exp.Consequence.Statements) != 1 {
        t.Errorf("Expected if consequence number of statements to be %d but got %d instead",
            1, len(exp.Consequence.Statements))
    }
    consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Consequence Statement is not an ExpressionStatement. Got %T instead", exp.Consequence.Statements[0])
    }
    testIdentifier(t, consequence.Expression, "x")

    // Check Alternative
    if len(exp.Alternative.Statements) != 1 {
        t.Errorf("Expected if consequence number of statements to be %d but got %d instead",
            1, len(exp.Alternative.Statements))
    }
    alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Alternative Statement is not an ExpressionStatement. Got %T instead", exp.Alternative.Statements[0])
    }
    testIdentifier(t, alternative.Expression, "y")
}

func TestFunctionLiteralParsing(t *testing.T) {
    input := "fn (x, y) { x + y; }"
    lex := lexer.NewLexer(input)
    par := NewParser(lex)
    pro := par.ParseProgram()

    if pro == nil {
        t.Fatalf("The program is nil")
    }
    if len(pro.Statements) != 1 {
        t.Errorf("Expected program to have %d statements but it have %d instead", 1, len(pro.Statements))
    }
    checkParserErrors(t, par)

    stm, ok := pro.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("The statement is not an ExpressionStatement. Got %T instead", pro.Statements[0])
    }
    lit, ok := stm.Expression.(*ast.FunctionLiteral)
    if !ok {
        t.Fatalf("The expression is not a FunctionLiteral. Got %t instead", stm.Expression)
    }

    // Check Function Parameters
    if len(lit.Parameters) != 2 {
        t.Errorf("Expected fuction literal to have %d parameters but got %d instead", 2, len(lit.Parameters))
    }
    testLiteralExpression(t, lit.Parameters[0], "x")
    testLiteralExpression(t, lit.Parameters[1], "y")

    // Check Function Body
    if len(lit.Body.Statements) != 1 {
        t.Errorf("Expected function literal body to have %d statements but got %d instead", 1, len(lit.Body.Statements))
    }
    bodyStm, ok := lit.Body.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Fuction Literal Body first statement is not an ExpressionStatement. Got %T instead",
            lit.Body.Statements[0])
    }
    testInfixExpression(t, bodyStm.Expression, "x", "+", "y")
}

func TestFunctionParametersParsing(t *testing.T) {
    tests := []struct {
          input string;      expectedParams []string
    } {
        { "fn () {}",        []string {} },
        { "fn (x, y) {}",    []string { "x", "y" } },
        { "fn (x, y, z) {}", []string { "x", "y", "z" } },
    }
    var acc bytes.Buffer
    for _, test := range tests {
        acc.WriteString(test.input + ";\n")
    }
    input := acc.String()
    lex := lexer.NewLexer(input)
    par := NewParser(lex)
    pro := par.ParseProgram()

    checkParserErrors(t, par)
    if pro == nil {
        t.Fatalf("Program is nil")
    }
    if len(pro.Statements) != len(tests) {
        t.Fatalf("Expected program number of statements to be %d but got %d instead", len(tests), len(pro.Statements))
    }

    for i, test := range tests {
        stm := pro.Statements[i].(*ast.ExpressionStatement)
        fnLit := stm.Expression.(*ast.FunctionLiteral)

        if len(fnLit.Parameters) != len(test.expectedParams) {
            t.Errorf("Expected Function Literal number of parameters to be %d but got %d instead",
                len(fnLit.Parameters), len(test.expectedParams))
        }
    }
}

func TestCallExpressionParsing(t *testing.T) {
    input := "add(1, 2 * 3, 4 + 5);"
    lex := lexer.NewLexer(input)
    par := NewParser(lex)
    pro := par.ParseProgram()

    if pro == nil {
        t.Fatalf("Program is nill")
    }
    if len(pro.Statements) != 1 {
        t.Fatalf("Expected program number of statements to be %d but got %d instead", 1, len(pro.Statements))
    }
    checkParserErrors(t, par)

    stm, ok := pro.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("The first statement is not an ExpressionStatement")
    }
    exp, ok := stm.Expression.(*ast.CallExpression)
    if !ok {
        t.Fatalf("The statement expression is not call expression")
    }

    // Check indentifier
    testIdentifier(t, exp.Function, "add")

    // Check arguments
    if len(exp.Arguments) != 3 {
        t.Errorf("Expected call expression to have %d argumenst but found %d instead", 3, len(exp.Arguments))
    }
    testLiteralExpression(t, exp.Arguments[0], 1)
    testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
    testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}
