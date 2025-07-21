// monkey/token/token.go

package token

const (
    // Special types
    ILLEGAL    = "ILLEGAL"
    EOF        = "EOF"

    // indentifiers + literals
    IDENT      = "IDENT" // add, foobar, x, y
    INT        = "INT"

    // Operators
    ASSIGN     = "="
    PLUS       = "+"
    MINUS      = "-"
    BANG       = "!"
    ASTERISK   = "*"
    SLASH      = "/"

    LT         = "<"
    GT         = ">"

    // Delimiters
    COMMA      = ","
    SEMICOLON  = ";"

    LPAREN     = "("
    RPAREN     = ")"
    LBRACE     = "{"
    RBRACE     = "}"

    // Keywords
    FUNCTION   = "FUNCTION"
    LET        = "LET"
    TRUE       = "TRUE"
    FALSE      = "FALSE"
    IF         = "IF"
    ELSE       = "ELSE"
    RETURN     = "RETURN"
)

type TokenType string

type Token struct {
    Type    TokenType
    Literal string
}

func NewToken(t TokenType, val byte) Token {
    return Token { Type: t, Literal: string(val) }
}

func NewTokenStr(t TokenType, val string) Token {
    return Token { Type: t, Literal: val }
}
