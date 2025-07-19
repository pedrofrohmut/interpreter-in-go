// monkey/lexer/lexer.go

package lexer

import (
	"monkey/token"
	"fmt"
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
	if lx.position >= (len(lx.input) - 1) {
		lx.ch = 0 // ascii for nul or eof
	} else {
		lx.ch = lx.input[lx.position + 1]
		lx.position += 1
		lx.readPosition += 1
	}
}

// Default constructor
func NewLexer(input string) *Lexer {
	// lx := &Lexer { input: input }
	// lx.nextChar()
	// return lx
	readPosition := 1
	ch := byte(0)
	if len(input) > 0 {
		ch = input[0]
	}
	if len(input) == 0 {
		readPosition = 0

	}
	return &Lexer {
		input: input,
		position: 0,
		readPosition: readPosition,
		ch: ch,
	}
}

func isIdentifierLetter(val byte) bool {
	if val >= 'a' && val <= 'z' { return true }
	if val >= 'A' && val <= 'Z' { return true }
	if val == '_' { return true }
	return false
}

func isNumber(val byte) bool {
	fmt.Println(val)
	if val >= '0' && val <= '9' {
		return true
	}
	return false
}

func (lx *Lexer) readIdentifierAsToken() token.Token {
	start := lx.position

	for lx.position < (len(lx.input) - 1) && isIdentifierLetter(lx.input[lx.position]) {
		lx.position += 1
		lx.readPosition += 1
	}
	identifier := lx.input[start:lx.position]

	fmt.Printf("IDENTIFIER: '%s'\n", identifier)

	switch identifier {
	case "let":
		return token.NewToken2(token.LET, identifier)
	default:
		return token.NewToken2(token.IDENT, identifier)
	}
}

// func (lx *Lexer) readIdentifierAsToken() token.Token {
// 	start := lx.position
//
// 	for isIdentifierLetter(lx.ch) {
// 		lx.nextChar()
// 	}
// 	identifier := lx.input[start:lx.position]
//
// 	// fmt.Printf("start: %q", start)
// 	// fmt.Printf("position: %q", lx.position)
// 	fmt.Printf("IDENTIFIER: '%s'\n", identifier)
//
// 	var tok token.Token
// 	switch identifier {
// 	case "let":
// 		tok = token.NewToken2(token.LET, identifier)
// 	default:
// 		tok = token.NewToken2(token.IDENT, identifier)
// 	}
//
// 	return tok
// }

func (lx *Lexer) readNumberAsToken() token.Token {
	start := lx.position

	for isNumber(lx.ch) {
		lx.nextChar()
	}
	num := lx.input[start:lx.position]

	return token.NewToken2(token.INT, num)
}

// Reads the current char the returns an equivalent token for it
// and interates the lexer to next char
func (lx *Lexer) NextToken() token.Token {
	if lx.ch == ' ' { fmt.Println("SPACE") }

	// Go to the next nonblank character from input
	for lx.ch == ' ' { lx.nextChar() }

	var tok token.Token

	switch lx.ch {
	// Operators
	// ASSIGN    = "="
	case '=':
		tok = token.NewToken(token.ASSIGN, lx.ch)
	// PLUS      = "+"
	case '+':
		tok = token.NewToken(token.PLUS, lx.ch)

	// LPAREN    = "("
	case '(':
		tok = token.NewToken(token.LPAREN, lx.ch)
	// RPAREN    = ")"
	case ')':
		tok = token.NewToken(token.RPAREN, lx.ch)
	// LBRACE    = "{"
	case '{':
		tok = token.NewToken(token.LBRACE, lx.ch)
	// RBRACE    = "}"
	case '}':
		tok = token.NewToken(token.RBRACE, lx.ch)

	// Delimiters
	// COMMA     = ","
	case ',':
		tok = token.NewToken(token.COMMA, lx.ch)
	// SEMICOLON = ";"
	case ';':
		tok = token.NewToken(token.SEMICOLON, lx.ch)

	case 0:
		tok = token.Token { Type: token.EOF, Literal: "" }

	default:
		if isIdentifierLetter(lx.ch) {
			tok = lx.readIdentifierAsToken()
		} else {
			tok = token.NewToken(token.ILLEGAL, lx.ch)
		}
	}

	return tok
}
