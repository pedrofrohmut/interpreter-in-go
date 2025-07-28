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
    fmt.Println("Monkey REPL. (:q to quit)")
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Printf(">> ")

        scanner.Scan()
        err := scanner.Err()
        if err != nil { log.Fatal(err) }

        line := scanner.Text()
        if line == ":q" { break }

        lx := lexer.NewLexer(line)
        lx.PrintTokens()
    }
}
