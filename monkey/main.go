// monkey/main.go

package main

import (
    "os"
    "fmt"
    "monkey/repl"
    "monkey/lexer"
    "monkey/parser"
    "monkey/evaluator"
)

func debugMain() {
    // Return Statement
    // input := "return a + b;"

    // Identifier Expressions
    // input := "foo;"

    // Precedence
    // input := "1 + 2 + 3;"
    // input := "a + b;"
    // input := "a + b + c;"
    // input := "a + b + c + d;"
    // input := "a + b; a + b + c; a + b + c + d;"
    // input := "a + b * c;"
    // input := "-a * b"
    // input := "a + (b * c);"

    // Prefix
    // input := "-1 + 2;"

    // IfElse
    // input := "if (foo < bar) { 13 } else { 42 };"
    // input := "if (foo < bar) { 13 };"

    // function literal
    // input := "fn (x, y) { x + y; }"
    // input := "fn (x, y) {}"
    // input := "fn () {}"

    // Call Expression
    // input := "add(1, 2 * 3, 4 + 5);"

    // Call Expression Precedence Test
    // input := "a + add(b * c) + d;" // => "((a + add((b * c))) + d)"
    // input := "add(b * c);" // => "add((b * c))"
    // input := "a + add(b * c);" // => "(a + add((b * c)))"
    // input := "add(b * c) + d;" // => "(add((b * c)) + d)"
    // input := "add();" // => "add()"

    // Return Statement
    // input := "return;" // => "return"
    // input := "return 5;" // => "return 5"
    // input := "return 5 + 10;" // => "return (5 + 10)"

    // Eval int
    // input := "5;"

    // Eval Bang
    // input := "!false;"

    // Eval Return
    // input := "9; return 3 * 7; 9;"
    // input := "return 3 * 7;"

    // Eval error
    input := "5 + true;"

    lexer := lexer.NewLexer(input)
    parser := parser.NewParser(lexer)
    program := parser.ParseProgram()
    // eva := evaluator.Eval(pro.Statements[0])
    eva := evaluator.Eval(program)

    // pro.PrintStatements()
    // pro.PrintStatements()
    // s := pro.String()
    // _ = s

    _ = eva
    _ = program
}

func replMain() {
    if len(os.Args) < 2 {
        fmt.Println("You did not provided the repl type. Add 'lexer', 'parser' or 'eval' as an argument.")
        os.Exit(0)
    }
    replType := os.Args[1]
    repl.Execute(replType)
}

func main() {
    const debug = true // Toggle for debugging or to use the repl
    if debug {
        debugMain()
    } else {
        replMain()
    }
}
