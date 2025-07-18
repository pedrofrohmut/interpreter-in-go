// monkey/token/token.go

package token

import (
	"fmt"
)

const (
	// Special types
	ILLEGAL   = "ILLEGAL"
	EOF       = "EOF"

	// indentifiers + literals
	IDENT     = "IDENT" // add, foobar, x, y
	INT       = "INT"

	// Operators
	ASSIGN    = "="
	PLUS      = "+"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	// Keywords
	FUNCTION  = "FUNCTION"
	LET       = "LET"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func newToken(t TokenType, val byte) Token {
	return Token { Type: t, Literal: string(val) }
}

func _newToken(t TokenType, val byte) *Token {
	return &Token { Type: t, Literal: string(val) }
}

func TokenFromValue(val byte) (Token, error) {
	switch val {
	// Operators
	// ASSIGN    = "="
	case '=':
		return newToken(ASSIGN, val), nil
	// PLUS      = "+"
	case '+':
		return newToken(PLUS, val), nil

	// LPAREN    = "("
	case '(':
		return newToken(LPAREN, val), nil
	// RPAREN    = ")"
	case ')':
		return newToken(RPAREN, val), nil
	// LBRACE    = "{"
	case '{':
		return newToken(LBRACE, val), nil
	// RBRACE    = "}"
	case '}':
		return newToken(RBRACE, val), nil

	// Delimiters
	// COMMA     = ","
	case ',':
		return newToken(COMMA, val), nil
	// SEMICOLON = ";"
	case ';':
		return newToken(SEMICOLON, val), nil

	case 0:
		return Token { Type: EOF, Literal: "" }, nil

	default:
		return Token {}, fmt.Errorf("Unexpected character. %c is not one of the expected characters on the list.", val)
	}
}
