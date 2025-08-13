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

    // Comparison
    LT         = "<"
    GT         = ">"
    EQ         = "=="
    NOT_EQ     = "!="

    // Delimiters
    COMMA      = ","
    SEMICOLON  = ";"

    // Grouping
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
    case PLUS, MINUS, BANG, ASTERISK, SLASH, LT, GT, EQ, NOT_EQ:
        return true
    default:
        return false
    }
}
