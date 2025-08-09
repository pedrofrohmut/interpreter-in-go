// monkey/main.go

package main

import (
    "os"
    "fmt"
    "monkey/repl"
    "monkey/lexer"
    "monkey/parser"
)

func debugMain() {
    // Identifier Expressions
    input := "foo;"
    // Precedence
    // input := "1 + 2 + 3;"

    // Prefix
    // input := "-1 + 2;"

    // IfElse
    // input := "if (foo < bar) { 13 } else { 42 }"

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

    lex := lexer.NewLexer(input)
    par := parser.NewParser(lex)
    pro := par.ParseProgram()
    _ = pro
    // s := pro.String()
    // _ = s
}

func replMain() {
    if len(os.Args) < 2 {
        fmt.Println("You did not provided the repl type. Add 'lexer' or 'parser' as an argument.")
        os.Exit(0)
    }
    replType := os.Args[1]
    repl.Execute(replType)
}

func main() {
    debugMain()

    // const debug = true // Toggle for debugging or to use the repl
    // if debug {
    //     debugMain()
    // } else {
    //     replMain()
    // }
}
