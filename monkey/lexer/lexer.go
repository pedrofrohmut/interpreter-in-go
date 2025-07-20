// monkey/lexer/lexer.go

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

func (lx *Lexer) getCh() byte {
	// Return 0 (ASCII character for EOF) when the position has reach the end of the input
	if lx.pos >= len(lx.input) {
		return 0
	}
	return lx.input[lx.pos]
}

func (lx *Lexer) nextPos() bool {
	if lx.pos < len(lx.input) {
		lx.pos += 1
		return true
	}
	return false
}

func (lx *Lexer) hasNextCh() bool {
	if lx.pos < len(lx.input) - 1 {
		return true
	}
	return false
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

func (lx *Lexer) skipWhiteSpaces() {
	for isWhiteSpace(lx.getCh()) {
		hasNext := lx.nextPos()
		if ! hasNext { break }
	}
}

func (lx *Lexer) readIdentifier() string {
	start := lx.pos
	for isIdentLetter(lx.input[lx.pos + 1]) {
		if ! lx.hasNextCh() { break }
		lx.pos += 1
	}
	return lx.input[start : lx.pos + 1]
}

// INFO: The monkey language, for educational purpose, only support int numbers.
// Other kinds of number can be later added as an exercise
func (lx *Lexer) readIntNumber() string {
	start := lx.pos
	for isIntNumber(lx.input[lx.pos + 1]) {
		if ! lx.hasNextCh() { break }
		lx.pos += 1
	}
	return lx.input[start:lx.pos + 1]
}

func (lx *Lexer) GetNextToken() token.Token {
	lx.skipWhiteSpaces()

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
	default:
		switch {
		case isIdentLetter(lx.getCh()) == true:
			ident := lx.readIdentifier()
			switch ident {
			case "let":
				tk = token.NewTokenStr(token.LET, ident)
			case "fn":
				tk = token.NewTokenStr(token.FUNCTION, ident)
			default:
				tk = token.NewTokenStr(token.IDENT, ident)
			}
		case isIntNumber(lx.getCh()) :
			num := lx.readIntNumber()
			tk = token.NewTokenStr(token.INT, num)
		default:
			tk = token.NewToken(token.ILLEGAL, lx.getCh())
		}
	}

	lx.nextPos()

	return tk
}

func (lx *Lexer) PrintChars() {
	start := lx.pos
	i := 0
	fmt.Printf("ASCII     \tDEC\tCHAR\n")
	for lx.hasNextCh() {
		fmt.Printf("Char[%d]: \t%d\t'%s'\n", i, lx.getCh(), string(lx.getCh()))
		lx.nextPos()
		i += 1
	}
	lx.pos = start
}
