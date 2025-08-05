// monkey/main.go

package main

import (
    "os"
    "monkey/repl"
    "monkey/lexer"
    "monkey/parser"
)


func main() {
    const debug = true
    if debug {
        // input := "1 + 2 + 3;" // 1st
        input := "-1 + 2" // 2nd
        lex := lexer.NewLexer(input)
        par := parser.NewParser(lex)
        _ = par.ParseProgram()
    } else {
        replType := os.Args[1]
        repl.Execute(replType)
    }
}
