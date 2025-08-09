// monkey/parser/parser.go

package parser

import (
    "fmt"
    "monkey/token"
    "monkey/lexer"
    "monkey/ast"
    "monkey/utils"
)

type Parser struct {
    lex *lexer.Lexer
    curr token.Token
    peek token.Token
    errors []string
}

func NewParser(lexer *lexer.Lexer) *Parser {
    parser := &Parser { lex: lexer }
    parser.curr = lexer.GetNextToken()
    parser.peek = lexer.GetNextToken()
    parser.errors = []string {}
    return parser
}

func (this *Parser) isCurr(tokenType string) bool {
    return this.curr.Type == tokenType
}

func (this *Parser) next() {
    this.curr = this.peek
    this.peek = this.lex.GetNextToken()
}

func (this *Parser) hasNext() bool {
    return this.curr.Type != token.EOF
}

func (this *Parser) addTokenError(tokenType string) {
    err := fmt.Sprintf("Expected token to be %s but got %s instead", tokenType, this.curr.Type)
    this.errors = append(this.errors, err)
}

func (this *Parser) parseLetStatement() *ast.LetStatement {
    if !this.isCurr(token.LET) { return nil }

    stm := ast.NewLetStatement()
    hasError := false

    // Check identifier
    this.next()
    if !this.isCurr(token.IDENT) {
        this.addTokenError(token.IDENT)
        hasError = true
    }
    if !hasError {
        stm.Identifier = this.curr.Literal
        this.next()
    }

    // Check for asign symbol
    if !this.isCurr(token.ASSIGN) {
        this.addTokenError(token.ASSIGN)
        hasError = true
    }
    this.next()

    // TODO: parse the expression later
    for !this.isCurr(token.SEMICOLON) { this.next() }

    if !this.isCurr(token.SEMICOLON) || hasError { return nil }

    return stm // Parse should end with curr == token.SEMICOLON
}

func (this *Parser) parseStatement() ast.Statement {
    return this.parseLetStatement()
}

func (this *Parser) ParseProgram() *ast.Program {
    program := ast.NewProgram()
    for this.hasNext() {
        stm := this.parseStatement()
        if !utils.IsNill(stm) {
            program.Statements = append(program.Statements, stm)
        }
        this.next()
    }
    return program
}
