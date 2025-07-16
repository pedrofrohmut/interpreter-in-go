package lexer

// The current version of Lexer only supports ASCII characters. Can be updated
// later to support utf-8 later as an exercise
type Lexer struct {
	input string
	position int     // current position in input (points to the current char)
	readPosition int // current read position in input (after current char)
	ch byte          // current char under examination
}

func NewLexer(initialInput string) *Lexer {
	newLexer := &Lexer{ input: initialInput }
	newLexer.readChar()
	return newLexer
}

func (lexer *Lexer) readChar() {
	if lexer.readPosition >= len(lexer.input) {
		lexer.ch = 0 // ascii for nul or eof
		return
	}
	lexer.ch = lexer.input[lexer.readPosition]
	lexer.position = lexer.readPosition
	lexer.readPosition += 1
}
