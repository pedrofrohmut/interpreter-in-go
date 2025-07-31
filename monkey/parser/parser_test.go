// monkey/parser/parser_test.go

package parser

import (
    "testing"
    _"reflect"
    _"fmt"
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

func TestLetStatement(t *testing.T) {
    input := `
        let x = 5;
        let y = 10;
        let z = 15;
    `
    lex := lexer.NewLexer(input)
    par := NewParser(lex)
    pro := par.ParseProgram()

    checkParserErrors(t, par)

    if pro == nil {
        t.Fatalf("Program is nill")
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
        if stm.Left.Value != test.expectedIdentifier {
            t.Errorf("[%d] Expected identifier to be %s but got %s", i, stm.Left.Value, test.expectedIdentifier)
        }
    }
}
// func TestLetStatement(t *testing.T) {
//     input := `
//         let x = 5;
//         let y = 10;
//         let z = 15;
//     `
//     lex := lexer.NewLexer(input)
//     par := NewParser(lex)
//     pro := par.ParseProgram()
//
//     if pro == nil {
//         t.Fatalf("Program is nill")
//     }
//     if len(pro.Statements) != 3 {
//         t.Fatalf("Expected program to have %d statements but got %d\n", 3, len(pro.Statements))
//     }
//     if len(par.errors) > 0 {
//         t.Fatalf("Parser the program found errors")
//     }
//
//     tests := []struct { expectedIdentifier string } {
//         {"x"}, {"y"}, {"z"},
//     }
//
//     for i, test := range tests {
//         stm, ok := pro.Statements[i].(*ast.LetStatement)
//         if !ok {
//             t.Errorf("Is not a LetStatement")
//         }
//         if stm.Token.Type != token.LET {
//             t.Errorf("[%dl] Expected identifier to be %s but got %s", i, token.LET, stm.Token.Type)
//         }
//         if stm.Identifier != test.expectedIdentifier {
//             t.Errorf("[%d] Expected identifier to be %s but got %s", i, stm.Identifier, test.expectedIdentifier)
//         }
//     }
// }
//
// func TestReturnStatement(t *testing.T) {
//     input := `
//         return 5;
//         return 10;
//         return 1234;
//     `
//     lex := lexer.NewLexer(input)
//     par := NewParser(lex)
//     pro := par.ParseProgram()
//
//     if pro == nil {
//         t.Fatalf("Program is nill")
//     }
//     if len(pro.Statements) != 3 {
//         t.Fatalf("Expected program to have %d statements but got %d\n", 3, len(pro.Statements))
//     }
//     if len(par.errors) > 0 {
//         t.Fatalf("Parser the program found errors")
//     }
//
//     for i, currStm := range pro.Statements {
//         stm, ok := currStm.(*ast.ReturnStatement)
//         if !ok {
//             t.Errorf("Is not a ReturnStatement")
//         }
//         if stm.Token.Type != token.RETURN {
//             t.Errorf("[%d] Expected token type to be %s but got %s", i, token.RETURN, stm.Token.Type)
//         }
//         if stm.Token.Literal != "return" {
//             t.Errorf("[%d] Expected token literal to be '%s' but got '%s'", i, "return", stm.Token.Literal)
//         }
//     }
// }
//
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
        t.Fatalf("Program is nill")
    }
    expectedErrCount := 4
    if len(par.errors) != expectedErrCount {
        t.Fatalf("Expected number of errors to be %d but got %d instead.", expectedErrCount, len(par.errors))
    }
}
// func TestErrorsOnLetStatement(t *testing.T) {
//     input := `
//         let x 5;
//         let = 10;
//         let 15;
//     `
//     lex := lexer.NewLexer(input)
//     par := NewParser(lex)
//     pro := par.ParseProgram()
//
//     // print("Program stm: ", len(pro.Statements), "\n")
//     // for i, stm := range pro.Statements {
//     //     fmt.Printf("[%d] ERR stm: %s\n", i, stm)
//     // }
//     // for i, err := range par.errors {
//     //     fmt.Printf("[%d] ERROR: %s\n", i, err)
//     // }
//
//     if pro == nil {
//         t.Fatalf("Program is nill")
//     }
//     expectedErrCount := 4
//     if len(par.errors) != expectedErrCount {
//         t.Fatalf("Expected number of errors to be %d but got %d instead.", expectedErrCount, len(par.errors))
//     }
// }
//
// func TestParserGetToken(t *testing.T) {
//     input := `
//         let x = 5;
//         let y = 10;
//     `
//     // input := ""
//     lex := lexer.NewLexer(input)
//     par := NewParser(lex)
//
//     tests := [] struct {
//         expectedType token.TokenType
//         expectedLiteral string
//     } {
//         // Statement 1
//         {token.LET, "let"},
//         {token.IDENT, "x"},
//         {token.ASSIGN, "="},
//         {token.INT, "5"},
//         {token.SEMICOLON, ";"},
//         // Statement 2
//         {token.LET, "let"},
//         {token.IDENT, "y"},
//         {token.ASSIGN, "="},
//         {token.INT, "10"},
//         {token.SEMICOLON, ";"},
//         // End
//         {token.EOF, ""},
//     }
//
//     // fmt.Printf("Current: %s\n\n", par.GetCurrToken())
//     for i, test := range tests {
//         tok := par.GetNextToken()
//         // fmt.Printf("[%d] %s\n",i, tok)
//         // fmt.Printf("Len: %d\n", len(par.tokens))
//         // fmt.Printf("Current: %s\n\n", par.GetCurrToken())
//         if tok.Type != test.expectedType {
//             t.Errorf("[%d] Expected token type to be %s but got %s instead", i, test.expectedType, tok.Type)
//         }
//         if tok.Literal != test.expectedLiteral {
//             t.Errorf("[%d] Expected token literal to be %s but got %s instead", i, test.expectedLiteral, tok.Literal)
//         }
//     }
// }
//
// func TestProgramString(t *testing.T) {
//     input := `
//         let x = 5;
//         let y = 10;
//     `
//     lex := lexer.NewLexer(input)
//     par := NewParser(lex)
//     pro := par.ParseProgram()
//
//     fmt.Printf("Program to string: '%s'\n", pro.String())
// }
//
// // func TestIdentifierExpression(t *testing.T) {
// //     input := "foo;"
// //     lex := lexer.NewLexer(input)
// //     par := NewParser(lex)
// //     pro := par.ParseProgram()
// //
// //     expectedLen := 1
// //     if len(pro.Statements) != expectedLen {
// //         t.Fatalf("Expected program statements length to be %d but got %d", expectedLen, len(pro.Statements))
// //     }
// //     expectedErrLen := 0
// //     if len(par.errors) != expectedErrLen {
// //         t.Errorf("Expected to not have any parser errors but got %d instead", len(par.errors))
// //     }
// //     stm, ok := pro.Statements[0].(*ast.ExpressionStatement)
// //     if !ok {
// //         t.Fatalf("Not an expression statement")
// //     }
// //     expectedExpression := "foo"
// //     if stm.Expression != expectedExpression {
// //         t.Errorf("Expected statement expression to be '%s' but got '%s' instead", expectedExpression, stm.Expression)
// //     }
// // }
//
// func TestIntExpression(t *testing.T) {
//     input := "5;"
//     lex := lexer.NewLexer(input)
//     par := NewParser(lex)
//     pro := par.ParseProgram()
//
//     expectedNumStm := 1
//     if len(pro.Statements) != expectedNumStm {
//         t.Fatalf("Expected program number of statements to be %d but got %d instead", expectedNumStm, len(pro.Statements))
//     }
//     stm, ok := pro.Statements[0].(*ast.ExpressionStatement)
//     if !ok {
//         t.Fatalf("Expected an ast.ExpressionStatement but got %T instead", pro.Statements[0])
//     }
//     // fmt.Println(stm)
//     if stm.Expression != "5" {
//         t.Fatalf("Expected statement expression to be '%s' but got '%s' instead", "5", stm.Expression)
//     }
// }
//
// func TestParsingPrefixExpression(t *testing.T) {
//     input := `
//         !5;
//         -15;
//     `
//     lex := lexer.NewLexer(input)
//     par := NewParser(lex)
//     pro := par.ParseProgram()
//
//     // for i, stm := range pro.Statements {
//     //     fmt.Printf("[%d] stm: %s\n", i, stm)
//     // }
//
//     // tests := []struct {
//     //     input string
//     //     operator string
//     //     integerValue int64
//     // } {
//     //     {"!5;", "!", 5},
//     //     {"-15;", "-", 15},
//     // }
//
//     tests := []struct { operator string; integerValue int64 } {
//         {"!", 5}, {"-", 15},
//     }
//
//     _, _ = tests, pro
//
//     if pro == nil {
//         t.Fatalf("Program is nill")
//     }
//     expectedNumStm := 2
//     if len(pro.Statements) != expectedNumStm {
//         t.Fatalf("Expected program number of statements to be %d but got %d instead", expectedNumStm, len(pro.Statements))
//     }
//
//     for i, err := range par.errors {
//         fmt.Printf("[%d] ERROR: %s\n", i, err)
//     }
//
//     for i, test := range tests {
//         stm, ok := pro.Statements[0].(*ast.ExpressionStatement)
//         if !ok {
//             t.Errorf("Not an expression statement")
//         }
//         fmt.Printf("[%d] Stm: %s\n", i, stm)
//         _, _ = stm, test
//     }
// }
