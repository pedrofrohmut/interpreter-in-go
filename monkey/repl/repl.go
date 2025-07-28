// monkey/repl/repl.go

package repl

import (
    "os"
    "bufio"
    "fmt"
    "log"
    "monkey/lexer"
)

func Execute() {
    fmt.Println("Monkey REPL.")
    fmt.Println("Tokenize your input")
    fmt.Println(":q or :quit to quit")
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
