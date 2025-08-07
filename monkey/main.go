// monkey/main.go

package main

import (
    "os"
    "monkey/repl"
    "monkey/lexer"
    "monkey/parser"
)

func debugMain() {
    // input := "1 + 2 + 3;" // precedence
    // input := "-1 + 2;" // prefix
    // input := "if (foo < bar) { 13 } else { 42 }" // ifelse
    // input := "fn (x, y) { x + y; }" // function literal
    // input := "fn (x, y) {}" // function literal
    input := "fn () {}" // function literal
    lex := lexer.NewLexer(input)
    par := parser.NewParser(lex)
    par.ParseProgram()
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
