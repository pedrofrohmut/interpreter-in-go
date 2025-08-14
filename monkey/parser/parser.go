// monkey/parser/parser.go

package parser

import (
    "fmt"
    "strconv"
    "monkey/ast"
    "monkey/lexer"
    "monkey/token"
    "monkey/utils"
)

const (
    // iota gives the constants a ascending numbers
    // _ skips the 0 value
    _ int = iota
    LOWEST
    EQUALS      // ==
    LESSGREATER // > or <
    SUM         // + -
    PRODUCT     // * /
    PREFIX      // -X or !X
    CALL        // myFunction(X)
    HIGHEST
)

var precedences = map[string] int {
    token.EQ:       EQUALS,
    token.NOT_EQ:   EQUALS,
    token.LT:       LESSGREATER,
    token.GT:       LESSGREATER,
    token.PLUS:     SUM,
    token.MINUS:    SUM,
    token.SLASH:    PRODUCT,
    token.ASTERISK: PRODUCT,
    token.LPAREN:   CALL,
}

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

func (this *Parser) isPeek(tokenType string) bool {
    return this.peek.Type == tokenType
}

func (this *Parser) next() {
    this.curr = this.peek
    this.peek = this.lex.GetNextToken()
}

func (this *Parser) hasNext() bool {
    return this.curr.Type != token.EOF
}

func (this *Parser) currPrecedence() int {
    return precedences[this.curr.Type]
}

func (this *Parser) peekPrecedence() int {
    return precedences[this.peek.Type]
}

func (this *Parser) addTokenError(tokenType string) {
    err := fmt.Sprintf("Expected token to be %s but got %s instead", tokenType, this.curr.Type)
    this.errors = append(this.errors, err)
}

func (this *Parser) addError(msg string) {
    this.errors = append(this.errors, msg)
}

func (this *Parser) parseLetStatement() *ast.LetStatement {
    if !this.isCurr(token.LET) { return nil } // BEFORE

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

    if hasError { return nil } // AFTER

    return stm // Parse should end with curr == token.SEMICOLON
}

func (this *Parser) parseReturnStatement() *ast.ReturnStatement {
    if !this.isCurr(token.RETURN) { return nil } // BEFORE

    stm := ast.NewReturnStatement()

    this.next()

    // TODO: parse the expression later
    for !this.isCurr(token.SEMICOLON) { this.next() }

    return stm
}

func (this *Parser) parsePrefixOrSymbol() ast.Expression {
    switch this.curr.Type {
    case token.BANG, token.MINUS:
        var pre = &ast.PrefixExpression {}
        pre.Operator = this.curr.Literal
        this.next()
        pre.Value = this.createNewInfixGroup(PREFIX)
        return pre
    case token.IDENT:
        return ast.NewIdentifier(this.curr.Literal)
    case token.INT:
        intValue, err := strconv.ParseInt(this.curr.Literal, 10, 64)
        if err != nil {
            this.addError("Could not convert current token literal to int64")
        }
        return ast.NewIntegerLiteral(intValue)
    default:
        this.addError("Invalid symbol or prefix to parse: " + this.curr.Type)
        return nil
    }
}

func (this *Parser) makeInfix(left ast.Expression) ast.Expression {
    // Curr is operator
    var inf = &ast.InfixExpression {}
    inf.Left = left
    inf.Operator = this.curr.Literal
    var precedence = this.currPrecedence()
    this.next() // Curr is next value
    inf.Right = this.createNewInfixGroup(precedence)
    return inf
}

func (this *Parser) createNewInfixGroup(ctxPrecedence int) ast.Expression {
    var parsedValue = this.parsePrefixOrSymbol()

    var acc = parsedValue

    for !this.isPeek(token.SEMICOLON) && this.peekPrecedence() > ctxPrecedence {
        this.next() // Curr is operator
        acc = this.makeInfix(acc)
    }

    return acc
}

/// 1. add to left every time precedence is the same or lower
/// 2. create a new group every time it goes up
/// 3. close the group when it goes lower again
func (this *Parser) parseExpression() ast.Expression {
    return this.createNewInfixGroup(LOWEST)
}

func (this *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    stm := &ast.ExpressionStatement {}
    stm.Expression = this.parseExpression()
    this.next()
    return stm
}

func (this *Parser) parseStatement() ast.Statement {
    switch this.curr.Type {
    case token.LET:
        return this.parseLetStatement()
    case token.RETURN:
        return this.parseReturnStatement()
    case token.ILLEGAL:
        return nil
    default:
        return this.parseExpressionStatement()
    }
}

func (this *Parser) ParseProgram() *ast.Program {
    program := ast.NewProgram()
    for this.hasNext() {
        stm := this.parseStatement()

        if utils.IsNill(stm) {
            this.addError("Statement is nill. Something went wrong!!!")
            return nil
        }

        if !this.isCurr(token.SEMICOLON) {
            this.addError("The statement did not end with a semicolon")
            return nil
        }

        program.Statements = append(program.Statements, stm)
        this.next() // Jumps the semicolon
    }
    return program
}
