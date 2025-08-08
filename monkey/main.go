// monkey/main.go

package main

import (
    "os"
    "monkey/repl"
    "monkey/lexer"
    "monkey/parser"
)

func debugMain() {
    // precedence
    // input := "1 + 2 + 3;"

    // prefix
    // input := "-1 + 2;"

    // ifelse
    // input := "if (foo < bar) { 13 } else { 42 }"

    // function literal
    // input := "fn (x, y) { x + y; }"

    // function literal
    // input := "fn (x, y) {}"

    // function literal
    // input := "fn () {}"

    // call expression
    // input := "add(1, 2 * 3, 4 + 5);"

    // Call Expression Precedence Test => "((a + add((b * c))) + d)"
    // input := "a + add(b * c) + d;"

    //  Call Expression Precedence Test => "add((b * c))"
    // input := "add(b * c);"

    //  Call Expression Precedence Test => "(a + add((b * c)))"
    // input := "a + add(b * c);"

    //  Call Expression Precedence Test => "(add((b * c)) + d)"
    // input := "add(b * c) + d;"

    //  Call Expression Precedence Test => "add()"
    input := "add();"

    lex := lexer.NewLexer(input)
    par := parser.NewParser(lex)
    pro := par.ParseProgram()
    s := pro.String()
    _ = s
}

func replMain() {
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
