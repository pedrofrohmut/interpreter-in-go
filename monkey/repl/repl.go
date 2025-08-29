// monkey/repl/repl.go

package repl

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "monkey/evaluator"
    "monkey/lexer"
    "monkey/object"
    "monkey/parser"
)

func lexerRepl() {
    fmt.Println("Tokenize your input")
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
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Printf(">> ")

        scanner.Scan()
        var err = scanner.Err()
        if err != nil { log.Fatal(err) }

        var line = scanner.Text()
        if line == ":q" || line == ":quit" { break }

        var lexer = lexer.NewLexer(line)
        var parser = parser.NewParser(lexer)
        var program = parser.ParseProgram()

        if len(parser.Errors()) > 0 {
            parser.PrintErrors()
        } else {
            program.PrintStatements()
        }
    }
}

func evalRepl() {
    fmt.Println("Tokenize then Parse and then Eval your input")
    scanner := bufio.NewScanner(os.Stdin)
    var env = object.NewEnvironment()
    for {
        fmt.Printf(">> ")

        scanner.Scan()
        var err = scanner.Err()
        if err != nil { log.Fatal(err) }

        var line = scanner.Text()
        if line == ":q" || line == ":quit" { break }

        var lexer = lexer.NewLexer(line)
        var parser = parser.NewParser(lexer)
        var program = parser.ParseProgram()

        if len(parser.Errors()) > 0 {
            parser.PrintErrors()
        }

        var obj = evaluator.Eval(program, env)
        if obj != nil {
            fmt.Println(obj.Inspect())
        } else {
            fmt.Println("WARN: Evaluation result is nil")
        }
    }
}

func Execute(replType string) {
    fmt.Println("Monkey REPL. [:q or :quit to quit]")
    switch replType {
    case "lexer":
        lexerRepl()
    case "parser":
        parserRepl()
    case "eval":
        evalRepl()
    default:
        fmt.Println("You need to pass what kind of REPL you want as argument.")
        fmt.Println("Options are: 'lexer' and 'parser'.")
    }
}
