// monkey/repl/repl.go

package repl

import (
    "io"
    "bufio"
    "fmt"
    "monkey/lexer"
    "monkey/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
    scanner := bufio.NewScanner(in)

    for {
        fmt.Printf(PROMPT)
        scanned := scanner.Scan()
        if !scanned { return }

        line := scanner.Text()
        lx := lexer.NewLexer(line)

        tk := lx.GetNextToken()
        for tk.Type != token.EOF {
            fmt.Printf("%+v\n", tk)
            tk = lx.GetNextToken()
        }
    }
}

