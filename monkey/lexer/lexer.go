// monkey/lexer/lexer.go
/*
    Turns the input into a array of tokens
*/

package lexer

import (
    "monkey/token"
    "fmt"
)

// INFO: The current version of Lexer only supports ASCII characters. Can be updated
// later to support utf-8 later as an exercise
type Lexer struct {
    input string
    pos int
}

func NewLexer(input string) *Lexer {
    return &Lexer { input: input, pos: 0 }
}

func (this *Lexer) getCh() byte {
    // Return 0 (ASCII character for EOF) when the position has reach the end of the input
    if this.pos >= len(this.input) {
        return 0
    }
    return this.input[this.pos]
}

func (this *Lexer) nextPos() bool {
    if this.pos < len(this.input) {
        this.pos += 1
        return true
    }
    return false
}

func (this *Lexer) hasNextCh() bool {
    if this.pos < len(this.input) - 1 {
        return true
    }
    return false
}

func (this *Lexer) getNextCh() byte {
    if this.hasNextCh() {
        return this.input[this.pos + 1]
    }
    return 0 // EOF
}

func isIdentLetter(val byte) bool {
    if val >= 'a' && val <= 'z' { // lowercase letters
        return true
    }
    if val >= 'A' && val <= 'Z' { // uppercase letters
        return true
    }
    if val == '_' { // allowed special characters
        return true
    }
    return false
}

func isIntNumber(val byte) bool {
    if val >= '0' && val <= '9' {
        return true
    }
    return false
}

func isWhiteSpace(val byte) bool {
    if val == '\t' || val == ' ' || val == '\n' || val == '\r'{
        return true
    }
    return false
}

func (this *Lexer) skipWhiteSpaces() {
    for isWhiteSpace(this.getCh()) {
        hasNext := this.nextPos()
        if ! hasNext { break }
    }
}

func (this *Lexer) readIdentifier() string {
    start := this.pos
    for this.hasNextCh() && isIdentLetter(this.input[this.pos + 1]) {
        this.pos += 1
    }
    return this.input[start : this.pos + 1]
}

// INFO: The monkey language, for educational purpose, only support int numbers.
// Other kinds of number can be later added as an exercise
func (this *Lexer) readIntNumber() string {
    start := this.pos
    for this.hasNextCh() && isIntNumber(this.input[this.pos + 1]) {
        this.pos += 1
    }
    return this.input[start:this.pos + 1]
}

func (this *Lexer) GetNextToken() token.Token {
    this.skipWhiteSpaces()

    var tk token.Token

    switch this.getCh() {
    // Operators & Comparison
    case '=':
        switch this.getNextCh() {
        case '=':
            tk = token.NewTokenStr(token.Eq, "==")
            this.nextPos() // Needed for 2 characters operators
        default:
            tk = token.NewToken(token.Assign, this.getCh())
        }
    case '+':
        tk = token.NewToken(token.Plus, this.getCh())
    case '-':
        tk = token.NewToken(token.Minus, this.getCh())
    case '!':
        switch this.getNextCh() {
        case '=':
            tk = token.NewTokenStr(token.NotEq, "!=")
            this.nextPos() // Needed for 2 characters operators
        default:
            tk = token.NewToken(token.Bang, this.getCh())
        }
    case '*':
        tk = token.NewToken(token.Asterisk, this.getCh())
    case '/':
        tk = token.NewToken(token.Slash, this.getCh())
    case '<':
        tk = token.NewToken(token.Lt, this.getCh())
    case '>':
        tk = token.NewToken(token.Gt, this.getCh())

    // Delimiter
    case ',':
        tk = token.NewToken(token.Comma, this.getCh())
    case ';':
        tk = token.NewToken(token.Semicolon, this.getCh())

    case '(':
        tk = token.NewToken(token.Lparen, this.getCh())
    case ')':
        tk = token.NewToken(token.Rparen, this.getCh())
    case '{':
        tk = token.NewToken(token.Lbrace, this.getCh())
    case '}':
        tk = token.NewToken(token.Rbrace, this.getCh())

    case 0:
        tk = token.NewTokenStr(token.Eof, "")
    default:
        switch {
        case isIdentLetter(this.getCh()) == true:
            ident := this.readIdentifier()
            switch ident {
            case "true":
                tk = token.NewTokenStr(token.True, ident)
            case "false":
                tk = token.NewTokenStr(token.False, ident)
            case "let":
                tk = token.NewTokenStr(token.Let, ident)
            case "fn":
                tk = token.NewTokenStr(token.Function, ident)
            case "return":
                tk = token.NewTokenStr(token.Return, ident)
            case "if":
                tk = token.NewTokenStr(token.If, ident)
            case "else":
                tk = token.NewTokenStr(token.Else, ident)
            default:
                tk = token.NewTokenStr(token.Ident, ident)
            }
        case isIntNumber(this.getCh()) :
            num := this.readIntNumber()
            tk = token.NewTokenStr(token.Int, num)
        default:
            tk = token.NewToken(token.Illegal, this.getCh())
        }
    }

    this.nextPos()

    return tk
}

func (this *Lexer) PrintChars() {
    start := this.pos
    i := 0
    fmt.Printf("ASCII     \tDEC\tCHAR\n")
    for this.hasNextCh() {
        fmt.Printf("Char[%d]: \t%d\t'%s'\n", i, this.getCh(), string(this.getCh()))
        this.nextPos()
        i += 1
    }
    this.pos = start
}

func (this *Lexer) PrintTokens() {
    i := 0
    tk := this.GetNextToken()
    for tk.Type != token.Eof {
        fmt.Printf("[%d] %+v\n", i, tk)
        tk = this.GetNextToken()
    }
}
