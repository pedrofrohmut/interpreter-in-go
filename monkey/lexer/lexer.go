// monkey/lexer/lexer.go

package lexer

import (
	"monkey/token"
)

// The current version of Lexer only supports ASCII characters. Can be updated
// later to support utf-8 later as an exercise
type Lexer struct {
	input string
	pos int
}

func NewLexer(input string) *Lexer {
	return &Lexer { input: input, pos: 0 }
}

func (lx *Lexer) getCh() byte {
	// Return 0 (ASCII character for EOF) when the position has reach the end of the input
	if lx.pos >= len(lx.input) {
		return 0
	}
	return lx.input[lx.pos]
}

func (lx *Lexer) nextPos() {
	if lx.pos < len(lx.input) {
		lx.pos += 1
	}
}

func (lx *Lexer) GetNextToken() token.Token {
	var tk token.Token

	switch lx.getCh() {
	case '=':
		tk = token.NewToken(token.ASSIGN, lx.getCh())
	case '+':
		tk = token.NewToken(token.PLUS, lx.getCh())
	case '(':
		tk = token.NewToken(token.LPAREN, lx.getCh())
	case ')':
		tk = token.NewToken(token.RPAREN, lx.getCh())
	case '{':
		tk = token.NewToken(token.LBRACE, lx.getCh())
	case '}':
		tk = token.NewToken(token.RBRACE, lx.getCh())
	case ',':
		tk = token.NewToken(token.COMMA, lx.getCh())
	case ';':
		tk = token.NewToken(token.SEMICOLON, lx.getCh())
	case 0:
		tk = token.NewTokenStr(token.EOF, "")
	}

	lx.nextPos()

	return tk
}
