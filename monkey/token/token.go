// token/token.go

package token

type TokenType string

type Token struct {
	Type TokenType
	Literal string
}

const (
	// Special types
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	// indentifiers + literals
	IDENT = "IDENT" // add, foobar, x, y
	INT = "INT"

	// Operators
	ASSIGN = "="
	PLUS = "+"

	// Delimiters
	COMMA = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET = "LET"
)

