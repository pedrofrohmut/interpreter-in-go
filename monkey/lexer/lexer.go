// monkey/lexer/lexer.go

package lexer

import (
	"monkey/token"
	"log"
)

// The current version of Lexer only supports ASCII characters. Can be updated
// later to support utf-8 later as an exercise
type Lexer struct {
	input string
	position int     // current position in input (points to the current char)
	readPosition int // current read position in input (after current char)
	ch byte          // current char under examination
}

// Reads the input position at readPosition and set its value at ch and update
// the position and readPosition. if readPosition is not valid set ch to 0
func (lx *Lexer) nextChar() {
	if lx.readPosition >= len(lx.input) {
		lx.ch = 0 // ascii for nul or eof
		return
	}
	lx.ch = lx.input[lx.readPosition]
	lx.position = lx.readPosition
	lx.readPosition += 1
}

func NewLexer(input string) *Lexer {
	lx := &Lexer { input: input }
	lx.nextChar()
	return lx
}

// Special types
// ILLEGAL   = "ILLEGAL"
// EOF       = "EOF"

// indentifiers + literals
// IDENT     = "IDENT" // add, foobar, x, y
// INT       = "INT"

// Keywords
// FUNCTION  = "FUNCTION"
// LET       = "LET"

// Reads the current char the returns an equivalent token for it
// and interates the lexer to next char
func (lx *Lexer) NextToken() token.Token {
	tok, err := token.TokenFromValue(lx.ch)
	if err != nil {
		log.Fatal(err)
	}
	lx.nextChar()
	return tok
}
