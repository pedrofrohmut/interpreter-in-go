// monkey/main.go

package main

import (
    "os"
    "monkey/repl"
    "monkey/lexer"
    "monkey/parser"
)

func debugMain() {
    // input := "1 + 2 + 3;" // 1st
    // input := "-1 + 2" // 2nd
    input := "if (foo < bar) { 13 } else { 42 }" // 3rd
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
