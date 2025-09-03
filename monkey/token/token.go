// monkey/token/token.go

package token

const (
    // Special types
    Illegal    = "ILLEGAL"
    Eof        = "EOF"

    // indentifiers + literals
    Ident      = "IDENT" // add, foobar, x, y
    Int        = "INT"

    // Types
    String     = "STRING"

    // Operators
    Assign     = "="
    Plus       = "+"
    Minus      = "-"
    Bang       = "!"
    Asterisk   = "*"
    Slash      = "/"

    // Comparison
    Lt         = "<"
    Gt         = ">"
    Eq         = "=="
    NotEq     = "!="

    // Delimiters
    Comma      = ","
    Semicolon  = ";"

    // Grouping
    Lparen     = "("
    Rparen     = ")"
    Lbrace     = "{"
    Rbrace     = "}"

    // Keywords
    Function   = "FUNCTION"
    Let        = "LET"
    True       = "TRUE"
    False      = "FALSE"
    If         = "IF"
    Else       = "ELSE"
    Return     = "RETURN"
)

type Token struct {
    Type string
    Literal string
}

func NewToken(tokenType string, value byte) Token {
    return Token { Type: tokenType, Literal: string(value) }
}

func NewTokenStr(tokenType string, value string) Token {
    return Token { Type: tokenType, Literal: value }
}

func IsOperator(token Token) bool {
    switch token.Type {
    case Plus, Minus, Bang, Asterisk, Slash, Lt, Gt, Eq, NotEq:
        return true
    default:
        return false
    }
}
