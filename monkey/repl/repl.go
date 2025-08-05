// monkey/repl/repl.go

package repl

import (
    "os"
    "bufio"
    "fmt"
    "log"
    "monkey/lexer"
    "monkey/parser"
)

func printQuitMsg() {
    fmt.Println(":q or :quit to quit")
}

func lexerRepl() {
    fmt.Println("Tokenize your input")
    printQuitMsg()
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Printf(">> ")

        scanner.Scan()
        err := scanner.Err()
        if err != nil { log.Fatal(err) }

        line := scanner.Text()
        if line == ":q" || line == ":quit" { break }

        lx := lexer.NewLexer(line)
        lx.PrintTokens()
    }
}

func parserRepl() {
    fmt.Println("Tokenize then Parse your input")
    printQuitMsg()
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Printf(">> ")

        scanner.Scan()
        err := scanner.Err()
        if err != nil { log.Fatal(err) }

        line := scanner.Text()
        if line == ":q" || line == ":quit" { break }

        lex := lexer.NewLexer(line)
        par := parser.NewParser(lex)
        pro := par.ParseProgram()

        fmt.Println("Program: ", pro.String())
        for i, stm := range pro.Statements {
            fmt.Printf("[%d]:\t%s\n", i, stm)
        }
    }
}

func Execute(replType string) {
    fmt.Println("Monkey REPL.")
    switch replType {
    case "lexer":
        lexerRepl()
    case "parser":
        parserRepl()
    default:
        fmt.Println("You need to pass what kind of REPL you want as argument.")
        fmt.Println("Options are: 'lexer' and 'parser'.")
    }
}
